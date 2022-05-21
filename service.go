package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type TxService interface {
	GetTxInfo(request GetTxInfoRequest) ([]StructEthResponse, error)
}

type txservice struct{}

func (txservice) GetTxInfo(request GetTxInfoRequest) ([]StructEthResponse, error) {
	// wallet := "0xdafc581f2938a2d586cac3992cc77f42661df743"
	// startBlock := "14793029"
	apiKey := APIKEY
	resp, _ := http.Get("https://api.etherscan.io/api?module=account&action=txlist&address=" + request.Wallet +
		"&startblock=" + request.StartBlock + "&endblock=latest&sort=asc&apikey=" + apiKey)
	if resp == nil {
		err := errors.New("Error")
		return nil, err
	}
	defer resp.Body.Close()
	Body, _ := ioutil.ReadAll(resp.Body)
	getTxInfoResult := new(GetTxInfoResult)
	json.Unmarshal(Body, getTxInfoResult)

	return getFilteredTx(getTxInfoResult), nil
}

func getFilteredTx(getTxInfoResult *GetTxInfoResult) []StructEthResponse {
	result := []StructEthResponse{}
	for _, txInfo := range getTxInfoResult.Result {
		ethValue, err := strconv.Atoi(txInfo.Value)
		if err != nil {
			fmt.Println(err)
		}
		if ethValue > 0 {
			result = append(result, txInfo)
		}
	}
	return result
}
