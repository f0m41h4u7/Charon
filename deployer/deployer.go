package main
import (
        "github.com/gin-gonic/gin"
        "github.com/docker/distribution/notifications"
        "fmt"
	"encoding/json"
//	"os"
)

//func sendPatch(img string) {
//	client := &http.Client{}
//
//	req, err := http.NewRequest("GET")
//
//	req, err := http.NewRequest("PATCH", os.Getenv("APISERVER"), )
//
//	req.Header.Add("Content-Type", "application/json-patch+json")
//	resp, err := client.Do(req)
//}

func rollout(c *gin.Context) {
	body := c.Request.Body
	decoder := json.NewDecoder(body)

	var envelope notifications.Envelope
	err := decoder.Decode(&envelope)
	if err != nil {
		fmt.Sprintf("Failed to decode envelope: %s\n", err)
		return
	}
	for index, event := range envelope.Events {
		fmt.Printf("Processing event %d of %d\n", index+1, len(envelope.Events))
		if event.Action == notifications.EventActionPush {
			img := event.Target.Repository + ":" + event.Target.Tag
			fmt.Println(img)
//			sendPatch(img)
		}
	}
	c.JSON(200, 0)
}

func rollback(c *gin.Context) {
	//body := c.Request.Body
	//decoder := json.NewDecoder(body)

	//sendPatch(img)
	c.JSON(200, 0)
}

func main() {
//	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/rollout", rollout)
	r.POST("/rollback", rollback)
	r.Run(":31337")
}
