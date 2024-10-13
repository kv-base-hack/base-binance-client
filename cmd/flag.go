package main

import (
	"time"

	"github.com/urfave/cli/v2"
)

const (
	binanceApiKes                 = "binance-api-keys"
	getPriceDurationFlag          = "get-price-duration"
	updateBinanceInfoDurationFlag = "update-binance-info-duration"
)

var commonFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    binanceApiKes,
		EnvVars: []string{"BINANCE_API_KEYS"},
	},
	&cli.DurationFlag{
		Name:    getPriceDurationFlag,
		Value:   time.Second,
		EnvVars: []string{"GET_PRICE_DURATION"},
	},
	&cli.DurationFlag{
		Name:    updateBinanceInfoDurationFlag,
		Value:   time.Second,
		EnvVars: []string{"UPDATE_BINANCE_INFO_DURATION"},
	},
}

func newAppFlags() (flags []cli.Flag) {
	return commonFlags
}
