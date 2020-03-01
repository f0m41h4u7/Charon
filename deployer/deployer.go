package main
import (
        "github.com/gin-gonic/gin"
        "github.com/docker/distribution/notifications"
        "fmt"
	"encoding/json"
	"k8s.io/client-go/util/retry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	corev1 "k8s.io/api/core/v1"
	"os"
	"net/http"
	"io/ioutil"
	"crypto/tls"
	"crypto/x509"
	"bytes"
)
type AppMetadata struct {
	name            string
        namespace       string
}
type AppSpec struct {
	image   string
}
type App struct {
	apiVersion      string
        kind            string
        metadata        interface{}
        spec            interface{}
}

// Send updates
func sendUpdate(name string, img string) {
	fmt.Println("Sending update to ", img)

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
			fmt.Printf("Failed to get latest version of Pod: %v. Creating new Pod.", getErr)

			// If doesn't exist, create new

			certPath := "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
			tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
			addr := "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + "/apis/app.custom.cr/v1alpha1/namespaces/default/apps"
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
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs: caCertPool,
					},
				},
			}

			newApp := App {
				apiVersion:	"app.custom.cr/v1alpha1",
				kind:		"App",
				metadata: AppMetadata {
					name:	name,
				},
				spec: AppSpec {
					image:	img,
				},
			}
			reqBody, jsonErr := json.Marshal(newApp)
			if jsonErr != nil {
				fmt.Println(jsonErr)
			}
			fmt.Println("Request body: ", reqBody)

			// Send request to create App
			req, err := http.NewRequest("POST", addr, bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", token)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return err
			}

			defer resp.Body.Close()
			// Print response
			fmt.Println(resp.Body)
			return nil
		}
		// If exists, send update
		updApp.Spec.Containers[0].Image = img
		_, updateErr := podClient.Update(updApp)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
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
		fmt.Sprintf("Failed to decode envelope: %s\n", err)
		return
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

// Handle rollback notifications
func rollback(c *gin.Context) {
	//body := c.Request.Body
	//decoder := json.NewDecoder(body)

	//sendUpdate(name, img)
	c.JSON(200, 0)
}

func main() {
//	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/rollout", rollout)
	r.POST("/rollback", rollback)
	r.Run(":31337")
}
