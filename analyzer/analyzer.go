package main

import (
	"context"
	"fmt"
	"log"
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

func getMetrics() ([]float64, string) {
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
	result, warnings, err := v1api.QueryRange(ctx, "go_goroutines{instance=\"test-app:1337\",job=\"test-app\"}", r)
	if err != nil {
		log.Fatal(fmt.Errorf("Error querying Prometheus: %w\n", err))
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	pairs := result.(model.Matrix)[0].Values
	var vals []float64
	for _, p := range pairs {
		vals = append(vals, float64(p.Value))
	}

	return vals, "test-app"
}

func anomalyDetect() {
	metrics, image := getMetrics()

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
	if probability >= 85.0 {
		alarm := Alarm {
			Image: image,
		}

		reqBody, err := json.Marshal(alarm)
		if err != nil {
			log.Fatal(err)
		}

		client := http.Client{}
		req, err := http.NewRequest("POST", "charon-deployer:31337/rollback", bytes.NewReader(reqBody))
		if err != nil {
			err = fmt.Errorf("Failed to send notification: %w", err)
			log.Fatal(err)
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	for {
		anomalyDetect()
	}
}
