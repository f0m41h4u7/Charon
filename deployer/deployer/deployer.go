package deployer

import (
	"crypto/x509"

	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var (
	certPath  = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	tokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	Address   string
)

type Deployer struct {
	Address    string
	token      string
	caCertPool *x509.CertPool
	podClient  v1.PodInterface
}

func NewDeployer() *Deployer {
	var d Deployer
	d.Address = Address
	d.setToken()
	d.setCertPool()
	d.createPodClient()
	return &d
}
