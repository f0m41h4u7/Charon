package deployer

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

	m "github.com/f0m41h4u7/Charon/models"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

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
		log.Fatal(fmt.Errorf("Failed to create in-cluster config: %w", err))
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create clientset: %w", err))
	}
	d.podClient = clientset.CoreV1().Pods(corev1.NamespaceDefault)
}

func (d *Deployer) createNewCR(name string, img string) {
	// Create updated json config for the App
	newApp := m.App{
		ApiVersion: "app.custom.cr/v1alpha1",
		Kind:       "App",
		Metadata: m.Metadata{
			Name: name,
		},
		Spec: m.Spec{
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
	req, err := http.NewRequest("POST", d.Address, bytes.NewReader(reqBody))
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

func (d *Deployer) SendPatch(name string, img string) {
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
	newApp := m.Patch{
		Spec: m.Spec{
			Image: img,
		},
	}
	reqBody, err := json.Marshal(newApp)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create cr spec: %v\n %w\n", newApp, err))
	}
	req, err := http.NewRequest("PATCH", d.Address, bytes.NewReader(reqBody))
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

func (d *Deployer) GetPreviousVersion(name string) string {
	registryAddr := "http://" + os.Getenv("REGISTRY") + "v2/" + name + "/tags/list"
	resp, err := http.Get(registryAddr)
	if err != nil {
		err = fmt.Errorf("Failed to get image tags: %w", err)
		log.Fatal(err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tl = m.TagsList{}
	err = json.Unmarshal(respBytes, &tl)
	if err != nil {
		err = fmt.Errorf("Failed to parse body: %s\n %w", resp.Body, err)
		log.Fatal(err)
	}
	return tl.Tags[1]
}
