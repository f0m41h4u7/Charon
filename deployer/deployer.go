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
)

// Send Patch request
func sendPatch(name string, img string) {
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
		result, getErr := podClient.Get(name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Pod: %v", getErr))
		}

		result.Spec.Containers[0].Image = img
		_, updateErr := podClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated pod...")
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
				sendPatch(event.Target.Repository, img)
			}
		}
	}
	c.JSON(200, 0)
}

// Handle rollback notifications
func rollback(c *gin.Context) {
	//body := c.Request.Body
	//decoder := json.NewDecoder(body)

	//sendPatch(img)
	c.JSON(200, 0)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/rollout", rollout)
	r.POST("/rollback", rollback)
	r.Run(":31337")
}
