package chia

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

func NewClient(CertificateFile string, PrivateKey string, CACertificatePath string) *ChiaClient {
	//load client ceritifcate ( found in .chia/mainnet/config/ssl)
	//~/.chia/mainnet/config/ssl/full_node/private_full_node.crt || ~/.chia/mainnet/config/ssl/full_node/private_full_node.key
	cert, err := tls.LoadX509KeyPair(CertificateFile, PrivateKey)
	if err != nil {
		log.Fatalf("Error occured while processing the specified certificates/keys \n: %v", err)
	}
	//load ca ceritifcate ( found in .chia/mainnet/config/ssl/ca/chia_ca.crt)
	caCert, err := ioutil.ReadFile(CACertificatePath)
	if err != nil {
		log.Fatalf("Error occured while processing the specified ca \n: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	c := &http.Client{Transport: transport}
	return &ChiaClient{
		HTTPClient: c,
	}
}
func (c *ChiaClient) GetChiaBlockchainState(url string) (ChiaBlockchainState, error) {
	req, _ := http.NewRequest("POST", url+"/"+"get_blockchain_state", nil)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured while processing connection to %v (get_blockchain_state)  \n: %v", url, err)
	}

	dump, err := httputil.DumpResponse(res, true)
	log.Println(string(dump))
	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse ChiaBlockchainState
	json.Unmarshal(responseBody, &ServiceResponse)
	log.Println(ServiceResponse)
	return ServiceResponse, nil
}
