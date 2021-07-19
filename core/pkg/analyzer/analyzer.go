package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lytics/anomalyzer"
)

var currentAnom = make(map[string]float64)

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

func average() float64 {
	var res float64
	for _, val := range currentAnom {
		res += val
	}
	return res / float64(len(currentAnom))
}

func anomalyDetect(metricName string) bool {
	query, image := queryMetric(metricName)
	var metrics []float64
	for _, q := range query {
		tmp := strings.Split(q[1].(string), `"`)
		mt, _ := strconv.ParseFloat(tmp[0], 64)
		metrics = append(metrics, mt)
	}

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
	currentAnom[metricName] = probability
	if average() > 0.85 {
		alarm := Alarm{
			Image: image,
		}

		reqBody, err := json.Marshal(alarm)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("ANOMALY! %f\n", probability)
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
		return true
	}
	return false
}

func RunAnalyzer() {
	metrics := getMetricNames()
	for {
		anom := false
		for _, mt := range metrics {
			anom = anomalyDetect(mt)
			if anom {
				break
			}
		}
		if anom {
			time.Sleep(5 * time.Minute)
		}
	}
}
