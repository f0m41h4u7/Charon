package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func getTime(delta time.Duration) string {
	time := time.Now().Add(delta).UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(time, 10)
}

type promValue struct {
	Timestamp float64
	Value     string
}
type promMetric struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}
type promResult struct {
	Metric promMetric  `json:"metric"`
	Value  []promValue `json:"value"`
}

type promResponse struct {
	Status string   `json:"status"`
	Data   promData `json:"data"`
}

type promData struct {
	ResultType string       `json:"resultType"`
	Result     []promResult `json:"result"`
}

func queryMetric(metricName string) ([]promValue, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	promHost := os.Getenv("PROMETHEUS_HOST")
	addr := promHost + "/api/v1/query_range?query=" + metricName + "&start=" + getTime(-100*time.Hour) + "&end=" + getTime(0) + "&step=" + time.Minute.String()
	resp, err := http.Get(addr)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to connect to Prometheus: %w\n", err))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to read response body (queryMetric): %w\n", err))
	}

	var mt promResponse
	err = json.Unmarshal(respBytes, &mt)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to parse json (queryMetric)\n Request: %s\n Body: %s\n Error: %w\n", addr, respBytes, err))
	}
	fmt.Printf("%v", mt.Data.Result)
	return mt.Data.Result[0].Value, mt.Data.Result[0].Metric.Instance
}

func getMetricNames() []string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	promHost := os.Getenv("PROMETHEUS_HOST")
	addr := promHost + "/api/v1/label/__name__/values"
	resp, err := http.Get(addr)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to connect to Prometheus: %w\n", err))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to read response body (getMetricNames): %w\n", err))
	}

	type response struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}

	var mt response
	err = json.Unmarshal(respBytes, &mt)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to parse json (getMetricNames): %w\n", err))
	}

	return mt.Data
}
