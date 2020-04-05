package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lytics/anomalyzer"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Alarm struct {
	Image string `json:"image"`
}

func MinMax(array []float64) (float64, float64) {
	var max float64 = array[0]
	var min float64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func getMetrics(metricName string) ([]float64, string) {
	client, err := api.NewClient(api.Config{
		Address: os.Getenv("PROMETHEUS_HOST"),
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Error creating client: %w\n", err))
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v1.Range{
		Start: time.Now().Add(-100 * time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	result, warnings, err := v1api.QueryRange(ctx, metricName, r)
	if err != nil {
		log.Fatal(fmt.Errorf("Error querying Prometheus: %w\n", err))
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	log.Printf("Connected to Prometheus... Querying metrics...\n")
	pairs := result.(model.Matrix)[0].Values
	var vals []float64
	for _, p := range pairs {
		vals = append(vals, float64(p.Value))
	}

	return vals, "test-app"
}

func anomalyDetect(metricName string) {
	metrics, image := getMetrics(metricName)

	min, max := MinMax(metrics)
	conf := &anomalyzer.AnomalyzerConf{
		Sensitivity: 0.01,
		UpperBound:  max + 1,
		LowerBound:  min - 1,
		ActiveSize:  1,
		NSeasons:    4,
		Methods:     []string{"diff", "fence", "highrank", "lowrank", "magnitude"},
	}

	anom, err := anomalyzer.NewAnomalyzer(conf, metrics)
	if err != nil {
		log.Fatal(err)
	}

	probability := anom.Eval()
	log.Printf("Metric: %s; Probability: %f\n", metricName, probability)
	if probability > 0.85 {
		alarm := Alarm{
			Image: image,
		}

		reqBody, err := json.Marshal(alarm)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("ANOMALY in metric %s: %f\n", metricName, probability)
		httpcli := &http.Client{}
		req, err := http.NewRequest("POST", "http://charon-deployer:31337/rollback", bytes.NewReader(reqBody))
		if err != nil {
			err = fmt.Errorf("Failed to send notification: %v\n%w", req, err)
			log.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := httpcli.Do(req)
		if err != nil {
			log.Fatal(fmt.Errorf("Failed to send notification; %w\n", err))
		}
		defer resp.Body.Close()
		log.Printf("Sent rollback request!")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	for {
		anomalyDetect("testMetrics0{instance=\"test-app:1337\",job=\"test-app\"}")
		anomalyDetect("testMetrics1{instance=\"test-app:1337\",job=\"test-app\"}")
		anomalyDetect("testMetrics2{instance=\"test-app:1337\",job=\"test-app\"}")
	}
}
