package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   TxService
}

func (mw loggingMiddleware) GetTxInfo(request GetTxInfoRequest) (output []StructEthResponse, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetTxInfo",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.GetTxInfo(request)
	return
}
