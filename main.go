package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)
type structEthResponse struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}
type structResult struct {
	Result []structEthResponse
}

func main() {
	wallet := "0xdafc581f2938a2d586cac3992cc77f42661df743"
	startBlock := "14793029"
	apiKey := "xxxxxxxxxxxxxxxx"
	resp, err := http.Get("https://api.etherscan.io/api?module=account&action=txlist&address=" + wallet +
		"&startblock=" + startBlock + "&endblock=latest&sort=asc&apikey=" + apiKey)

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	Body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	ethResponse := new(structResult)
	json.Unmarshal(Body, ethResponse)


	for _, element := range ethResponse.Result {
		fmt.Printf("%+v\n", element)
	}

}
