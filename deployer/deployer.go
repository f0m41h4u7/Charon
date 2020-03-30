package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/docker/distribution/notifications"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
)

type AppMetadata struct {
	Name string `json:"name"`
}
type AppSpec struct {
	Image string `json:"image"`
}
type App struct {
	ApiVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Metadata   interface{} `json:"metadata"`
	Spec       interface{} `json:"spec"`
}
type Patch struct {
	Spec interface{} `json:"spec"`
}

// Send updates
func sendUpdate(name string, img string) {
	fmt.Println("Sending update to ", img)
	img = "charon-registry:5000/" + img

	// Get authorization token and certificate
	certPath := "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	addr := "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + "/apis/app.custom.cr/v1alpha1/namespaces/default/apps/" + name
	read, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		fmt.Println("Cannot read token", err)
	}
	token := "Bearer " + string(read)

	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		fmt.Println("Cannot get cert")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create HTTP client
	httpcli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	podClient := clientset.CoreV1().Pods(corev1.NamespaceDefault)

	// Try to update
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		updApp, getErr := podClient.Get(name, metav1.GetOptions{})
		if getErr != nil {
			fmt.Printf("Failed to get latest version of Pod: %v. Creating new Pod. \n", getErr)

			// Create updated json config for the App
			newApp := App{
				ApiVersion: "app.custom.cr/v1alpha1",
				Kind:       "App",
				Metadata: AppMetadata{
					Name: name,
				},
				Spec: AppSpec{
					Image: img,
				},
			}

			reqBody, err := json.Marshal(newApp)
			if err != nil {
				log.Fatal(fmt.Errorf("Failed to create cr spec: %v\n %w\n", newApp, err))
			}

			// Send request to create App
			req, err := http.NewRequest("POST", addr, bytes.NewReader(reqBody))
			if err != nil {
				log.Fatal(fmt.Errorf("Failed to send create request: %w\n", err))
			}
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", token)
			resp, err := httpcli.Do(req)
			if err != nil {
				return fmt.Errorf("Failed to create cr; %w\n", err)
			}

			defer resp.Body.Close()
			return nil
		}

		// If exists, send patch to app cr
		newApp := Patch{
			Spec: AppSpec{
				Image: img,
			},
		}
		reqBody, err := json.Marshal(newApp)
		if err != nil {
			log.Fatal(fmt.Errorf("Failed to create cr spec: %v\n %w\n", newApp, err))
		}
		req, err := http.NewRequest("PATCH", addr, bytes.NewReader(reqBody))
		if err != nil {
			log.Fatal(fmt.Errorf("Failed to send patch; %w\n", err))
		}

		req.Header.Add("Content-Type", "application/merge-patch+json")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", token)
		resp, err := httpcli.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}

		defer resp.Body.Close()

		// Update pod
		updApp.Spec.Containers[0].Image = img
		_, updateErr := podClient.Update(updApp)
		return updateErr
	})
	if retryErr != nil {
		log.Fatal(fmt.Errorf("Update failed: %v", retryErr))
	} else {
		fmt.Println("Successfully updated pod.")
	}
}

// Handle Registry notifications
func rollout(c *gin.Context) {
	body := c.Request.Body
	decoder := json.NewDecoder(body)

	// Receive envelope with evenets and decode it
	var envelope notifications.Envelope
	err := decoder.Decode(&envelope)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to decode envelope: %s\n", err))
	}

	// Process events
	for index, event := range envelope.Events {
		if event.Action == notifications.EventActionPush {
			fmt.Printf("Processing event %d of %d\n", index+1, len(envelope.Events))
			if (event.Target.Tag != "") && (event.Target.Repository != "charon-operator") && (event.Target.Repository != "deployer") {
				img := event.Target.Repository + ":" + event.Target.Tag
				fmt.Println(img)
				sendUpdate(event.Target.Repository, img)
			}
		}
	}
	c.JSON(200, 0)
}

type Anomaly struct {
	Image string `json:"image"`
}

// Handle rollback notifications
func rollback(c *gin.Context) {
	body := c.Request.Body
	var anom = Anomaly{}
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &anom)
	if err != nil {
		err = fmt.Errorf("Failed to parse body: %s\n %w", body, err)
		log.Fatal(err)
	}
	

	//sendUpdate(name, img)
	c.JSON(200, 0)
}

func main() {
	//	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/rollout", rollout)
	r.POST("/rollback", rollback)
	err := r.Run(":31337")
	if err != nil {
		log.Fatal(err)
	}
}
