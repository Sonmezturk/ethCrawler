package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func getTxInfoEndpoint(svc TxService) endpoint.Endpoint {

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
