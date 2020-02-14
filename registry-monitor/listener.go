package main
import (
        "github.com/gin-gonic/gin"
	"github.com/docker/distribution/notifications"
	"fmt"
//	"io/ioutil"
	"encoding/json"
)

func parseEvents(c *gin.Context) {
	body := c.Request.Body
//	x, _ := ioutil.ReadAll(body)
//	fmt.Printf("%s \n", string(x))

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
			fmt.Println("repositoryname: ", event.Target.Repository)
		}
	}
	c.JSON(200, 0)
}

func main() {
//	gin.SetMode(gin.ReleaseMode)
        r := gin.Default()
//	r.Use(parseEvents)
	r.POST("/event", parseEvents)
        r.Run(":31337")
}
