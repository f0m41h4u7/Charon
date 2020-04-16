package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

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
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type promData struct {
	ResultType string       `json:"resultType"`
	Result     []promResult `json:"result"`
}

//curl -g 'http://167.172.137.177:30329/api/v1/query?' --data-urlencode 'query=scrape_duration_seconds'
//{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"scrape_duration_seconds","instance":"test-app:1337","job":"test-app"},"value":[1586984270.136,"5.000423151"]}]}}

func queryMetric(metricName string) ([]promValue, string) {
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
	err = json.Unmarshal(respBytes, &mt)
	if err != nil {
		log.Fatal(err)
	}
	return mt.Data.(promData).Result[0].Value, mt.Data.(promData).Result[0].Metric.Instance
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
	err = json.Unmarshal(respBytes, &mt)
	if err != nil {
		log.Fatal(err)
	}

	return mt.Data.([]string)
}
