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
	"strings"

	"github.com/docker/distribution/notifications"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

var (
	certPath  = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	tokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	address   string
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

type Deployer struct {
	address    string
	token      string
	caCertPool *x509.CertPool
	podClient  v1.PodInterface
}

func (d *Deployer) setToken() {
	read, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Fatal(fmt.Errorf("Cannot read token, %w\n", err))
	}
	d.token = "Bearer " + string(read)
}

func (d *Deployer) setCertPool() {
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Fatal(fmt.Errorf("Cannot get cert, %w\n", err))
	}
	d.caCertPool = x509.NewCertPool()
	d.caCertPool.AppendCertsFromPEM(caCert)
}

func (d *Deployer) createPodClient() {
	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create in-cluster config: %w", err.Error()))
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create clientset: %w", err.Error()))
	}
	d.podClient = clientset.CoreV1().Pods(corev1.NamespaceDefault)
}

func newDeployer() *Deployer {
	var d Deployer
	d.address = address
	d.setToken()
	d.setCertPool()
	d.createPodClient()
	return &d
}

func (d *Deployer) createNewCR(name string, img string) {
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

	// Create HTTP client
	httpcli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: d.caCertPool,
			},
		},
	}

	reqBody, err := json.Marshal(newApp)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create cr spec: %v\n %w\n", newApp, err))
	}

	// Send request to create App
	req, err := http.NewRequest("POST", d.address, bytes.NewReader(reqBody))
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to send create request: %w\n", err))
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", d.token)
	resp, err := httpcli.Do(req)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create cr; %w\n", err))
	}
	defer resp.Body.Close()
}

func (d *Deployer) sendPatch(name string, img string) {
	registryName := os.Getenv("REGISTRY")
	img = registryName + img
	updApp, err := d.podClient.Get(name, metav1.GetOptions{})
	if err != nil {
		d.createNewCR(name, img)
		fmt.Println("Created new CR")
		return
	}
	// Create HTTP client
	httpcli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: d.caCertPool,
			},
		},
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
	req, err := http.NewRequest("PATCH", d.address, bytes.NewReader(reqBody))
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to send patch; %w\n", err))
	}

	req.Header.Add("Content-Type", "application/merge-patch+json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", d.token)
	resp, err := httpcli.Do(req)
	if err != nil {
		log.Fatal(fmt.Errorf("Patch failed: %v", err))
	}
	defer resp.Body.Close()

	// Update pod
	updApp.Spec.Containers[0].Image = img
	_, updErr := d.podClient.Update(updApp)
	if updErr != nil {
		log.Fatal(fmt.Errorf("Update failed: %v", updErr))
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
			name := event.Target.Repository
			if (event.Target.Tag != "") && (!strings.Contains(name, "charon-")) {
				img := name + ":" + event.Target.Tag
				fmt.Println(img)
				d := newDeployer()
				d.address = address + name
				d.sendPatch(name, img)
			}
		}
	}
	c.JSON(200, 0)
}

type Alarm struct {
	Image string `json:"image"`
}

// Handle rollback notifications
func rollback(c *gin.Context) {
	body := c.Request.Body
	var anom = Alarm{}
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &anom)
	if err != nil {
		err = fmt.Errorf("Failed to parse body: %s\n %w", body, err)
		log.Fatal(err)
	}

	d := newDeployer()
	// Find out, which version is deployed

	imgName := os.Getenv("REGISTRY") + anom.Image
	d.sendPatch(imgName, imgName)
	c.JSON(200, 0)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	address = "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + "/apis/app.custom.cr/v1alpha1/namespaces/default/apps/"

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/rollout", rollout)
	r.POST("/rollback", rollback)
	err = r.Run(":31337")
	if err != nil {
		log.Fatal(err)
	}
}
