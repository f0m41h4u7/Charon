package main

import (
	"log"
	"os"

	d "github.com/f0m41h4u7/Charon/deployer"
	"github.com/gin-gonic/gin"
)

func main() {
	d.Address = "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + "/apis/app.custom.cr/v1alpha1/namespaces/default/apps/"

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/rollout", rollout)
	r.POST("/rollback", rollback)
	err := r.Run(":31337")
	if err != nil {
		log.Fatal(err)
	}
}
