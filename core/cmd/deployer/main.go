package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/docker/distribution/notifications"
	"github.com/f0m41h4u7/Charon/core/pkg/deployer"
	"github.com/gin-gonic/gin"
)

// Alarm received from Analyzer, providing name of image with anomaly
type alarm struct {
	Image string `json:"image"`
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
			name := event.Target.Repository
			if event.Target.Tag != "" {
				img := name + ":" + event.Target.Tag
				fmt.Println(img)
				d := deployer.NewDeployer()
				d.Address += name
				d.SendPatch(name, img)
			}
		}
	}
	c.JSON(200, 0)
}

// Handle Analyzer notifications
func rollback(c *gin.Context) {
	body := c.Request.Body
	var anom = alarm{}
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &anom)
	if err != nil {
		err = fmt.Errorf("Failed to parse body: %s\n %w", body, err)
		log.Fatal(err)
	}
	log.Printf("Received a rollback request for image: %s\n", anom.Image)

	d := deployer.NewDeployer()
	// Find out, which version is deployed
	targetVersion := d.GetPreviousVersion(anom.Image)

	// Patch CR and rollback pod
	d.SendPatch(anom.Image, anom.Image+":"+targetVersion)
	c.JSON(200, 0)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/rollout", rollout)
	r.POST("/rollback", rollback)
	err := r.Run(":31337")
	if err != nil {
		log.Fatal(err)
	}
}
