package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type promMetric struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}
type promResult struct {
	Metric promMetric    `json:"metric"`
	Value  []interface{} `json:"value"`
}

type promResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type promData struct {
	ResultType string       `json:"resultType"`
	Result     []promResult `json:"result"`
}

//curl -g 'http://167.172.137.177:30329/api/v1/query?' --data-urlencode 'query=scrape_duration_seconds'
//{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"scrape_duration_seconds","instance":"test-app:1337","job":"test-app"},"value":[1586984270.136,"5.000423151"]}]}}

func queryMetric(metricName string) {
	promHost := os.Getenv("PROMETHEUS_HOST")
	resp, err := http.Get(promHost + "/api/v1/query_range?query=" + metricName + "&start=" + time.Now().Add(-100*time.Hour).String() + "&end=" + time.Now().String() + "&step=" + time.Minute.String())
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to connect to Prometheus: %w\n", err))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var mt promResponse
	err = json.Unmarshal(respBytes, mt)
	if err != nil {
		log.Fatal(err)
	}

}

func getMetrics(metricName string) []float64 {
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

	return vals
}

func getMetricNames() []string {
	promHost := os.Getenv("PROMETHEUS_HOST")
	resp, err := http.Get(promHost + "/api/v1/label/__name__/values")
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to connect to Prometheus: %w\n", err))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var mt promResponse
	err = json.Unmarshal(respBytes, mt)
	if err != nil {
		log.Fatal(err)
	}

	return mt.Data.([]string)
}
