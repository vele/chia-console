package chia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func FetchChiaMap() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/map", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "c7d2a342-a487-48b4-992e-27244df26506")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	var ServiceResponse ChiaCryptocurrencyMapResponse
	json.Unmarshal(responseBody, &ServiceResponse)
	bar := progressbar.Default(int64(len(ServiceResponse.Data)))

	for i := range ServiceResponse.Data {
		//fmt.Println(i, ServiceResponse.Data[i].Name)
		err := TableChiaMapInsert(i, ServiceResponse.Data[i].ID, ServiceResponse.Data[i].Rank, ServiceResponse.Data[i].Name, ServiceResponse.Data[i].Symbol, time.Now().String())
		if err != nil {
			return fmt.Errorf("Error when processing query: %v", err)
		}
		bar.Add(1)

	}
	return nil
}
func FetchChiaPrice(coin_id int) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	q := url.Values{}
	q.Add("id", fmt.Sprintf("%d", coin_id))

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "c7d2a342-a487-48b4-992e-27244df26506")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	var ServiceResponse CoinMarketCapSymbolResponse
	json.Unmarshal(responseBody, &ServiceResponse)
	err = TableChiaPricesInsert(&ServiceResponse)
	if err != nil {
		return fmt.Errorf("An error has occured: %s", err)
	}
	return nil
}
