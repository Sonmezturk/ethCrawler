package main

import (
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
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

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	var svc TxService
	svc = txservice{}
	svc = loggingMiddleware{logger, svc}

	getTxInfo := httptransport.NewServer(
		getTxInfoEndpoint(svc),
		decodeRequest,
		encodeResponse,
	)
	mux := http.NewServeMux()
	mux.Handle("/txInfo", getTxInfo)
	http.ListenAndServe(":8080", mux)

}
