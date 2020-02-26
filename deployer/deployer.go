/*****************************************************************************************
*
* The Deployer monitors notifications from Charon-registry and Charon-AI.
* Rollout/rollback updates are deployed by sending Patch requests to Kubernetes Apiserver.
*
******************************************************************************************/

package main
import (
        "github.com/gin-gonic/gin"
        "github.com/docker/distribution/notifications"
        "fmt"
	"encoding/json"
	"net/http"
	"bytes"
)

type metadataStruct struct {
	name string
}

type specStruct struct {
	image string
}

// Structure of deployer cr.yaml
type appConfig struct {
	apiVersion string
	kind string
	metadata metadataStruct
	spec specStruct
}

// Send Patch request to Apiserver
func sendPatch(img string) {
	client := &http.Client{}
	fmt.Println("Sending patch...")

	// Custom Resource config
	config := appConfig {
		apiVersion: "app.custom.cr/v1alpha1",
		kind: "App",
		metadata: metadataStruct {
			name: img + "-pod",
		},
		spec: specStruct {
			image: img,
		},
	}

	var patch []byte
	patch, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
	}

	addr := "http://10.96.0.1"
	req, err := http.NewRequest("PATCH", addr, bytes.NewReader(patch))
	req.Header.Add("Content-Type", "application/json-patch+json")
        resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending patch request: ", err)
	}
	fmt.Println(resp)
}

// Handle Charon-registry notifications
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
				sendPatch(img)
			}
		}
	}
	c.JSON(200, 0)
}

// Handle Charon-AI notifications
func rollback(c *gin.Context) {
	//body := c.Request.Body
	//decoder := json.NewDecoder(body)

	//sendPatch(img)
	c.JSON(200, 0)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Listen Charon-registry notifications
	r.POST("/rollout", rollout)

	// Listen Charon-AI notifications
	r.POST("/rollback", rollback)
	r.Run(":31337")
}
