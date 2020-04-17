package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type promResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name     string `json:"__name__"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func getTime(delta time.Duration) string {
	return strconv.FormatInt(time.Now().Add(delta).Unix(), 10)
}

func queryMetric(metricName string) ([][]interface{}, string) {
	promHost := os.Getenv("PROMETHEUS_HOST")
	params := metricName + "&start=" + getTime(-100*time.Hour) + "&end=" + getTime(0) + "&step=60s"
	paramsUrlEncoded := &url.URL{Path: params}
	addr := promHost + "/api/v1/query_range?query=" + paramsUrlEncoded.String()

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

	return mt.Data.Result[0].Values, mt.Data.Result[0].Metric.Instance
}

func getMetricNames() []string {
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
