package chia

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	requestBody, _ := json.Marshal(map[string]string{})
	req, _ := http.NewRequest("POST", url+"/"+"get_blockchain_state", bytes.NewBuffer(requestBody))
	req.Header.Add("Content-Type", "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured while processing connection to %v (get_blockchain_state)  \n: %v", url, err)
	}

	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse ChiaBlockchainState
	json.Unmarshal(responseBody, &ServiceResponse)
	log.Println(ServiceResponse)
	return ServiceResponse, nil
}
func (c *ChiaClient) GetChiaWallet(url string) (WalletBallance, error) {
	requestBody, _ := json.Marshal(map[string]string{"wallet_id": "1"})
	req, _ := http.NewRequest("POST", url+"/"+"get_wallet_balance", bytes.NewBuffer(requestBody))
	req.Header.Add("Content-Type", "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured while processing connection to %v (get_wallet_balance)  \n: %v", url, err)
	}

	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse WalletBallance
	json.Unmarshal(responseBody, &ServiceResponse)
	log.Println(ServiceResponse)
	return ServiceResponse, nil
}
func (c *ChiaClient) GetChiaPlots(url string) (ChiaPlots, error) {
	requestBody, _ := json.Marshal(map[string]string{"wallet_id": "1"})
	req, _ := http.NewRequest("POST", url+"/"+"get_plots", bytes.NewBuffer(requestBody))
	req.Header.Add("Content-Type", "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured while processing connection to %v (get_wallet_balance)  \n: %v", url, err)
	}

	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse ChiaPlots
	json.Unmarshal(responseBody, &ServiceResponse)
	return ServiceResponse, nil
}
func ParseLogs(logFile string) []string {
	fmt.Println(logFile)
	f, err := os.Open("/root/.chia/mainnet/log/debug.log")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	log := &Line{}
	var logs []string
	for sc.Scan() {
		ok, err := log.Extract(sc.Bytes())
		if !ok {
			if err != nil {
				fmt.Println(err)
			}
			logs = append(logs, string(log.Plots))
			continue
		}
	}
	return logs
}
