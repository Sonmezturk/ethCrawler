package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/schema"
)
var APIKEY = os.Getenv("APIKEY")
var decoder = schema.NewDecoder()

type StructEthResponse struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type GetTxInfoResult struct {
	Result []StructEthResponse `json:"result"`
}

type GetTxInfoRequest struct {
	Wallet     string `schema:"wallet"`
	StartBlock string `schema:"startBlock"`
}

type Service interface {
	GetTxInfo(request GetTxInfoRequest) ([]StructEthResponse, error)
}

type service struct{}

func (service) GetTxInfo(request GetTxInfoRequest) ([]StructEthResponse, error) {
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

func getTxInfoEndpoint(svc Service) endpoint.Endpoint {

	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(GetTxInfoRequest)
		res, err := svc.GetTxInfo(req)
		if err != nil {
			return nil, errors.New("Error")
		}
		return GetTxInfoResult{res}.Result, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request GetTxInfoRequest
	err := decoder.Decode(&request, r.URL.Query())
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {
	svc := service{}

	getTxInfo := httptransport.NewServer(
		getTxInfoEndpoint(svc),
		decodeRequest,
		encodeResponse,
	)
	mux := http.NewServeMux()
	mux.Handle("/txInfo", getTxInfo)
	log.Fatal(http.ListenAndServe(":8080", mux))

}
