package bf

import (
	"net"
	"net/http"
	"time"

	"github.com/adshao/go-binance/v2"
	bl "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

const timeout = time.Second * 5

func getBinanceClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          64,
		MaxIdleConnsPerHost:   16,
		MaxConnsPerHost:       64,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: time.Second * 30,
	}

	return &http.Client{Transport: NewTransportRateLimiter(transport)}
}

func NewBinance(apiKey string, secret string) (*bl.Client, *futures.Client) {
	spot := binance.NewClient(apiKey, secret)
	bc := bl.NewFuturesClient(apiKey, secret)
	bc.HTTPClient = getBinanceClient()
	spot.HTTPClient = getBinanceClient()
	return spot, bc
}
