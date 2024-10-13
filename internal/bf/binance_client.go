package bf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kv-base-hack/kv-binance/common"
)

const (
	apiBaseURL   = "https://api.binance.com"
	apiKeyHeader = "X-MBX-APIKEY"
)

// Client to interact with binance api
type Client struct {
	httpClient *http.Client
	apiKey     string
	secretKey  string
}

// FwdData contain data we forward to client
type FwdData struct {
	Status      int
	ContentType string
	Data        []byte
}

// NewClient create new client object
func NewClient(key, secret string, client *http.Client) *Client {
	return &Client{
		httpClient: client,
		apiKey:     key,
		secretKey:  secret,
	}
}

// ListenKey is listen for user data stream
type ListenKey struct {
	ListenKey string `json:"listenKey"`
}

func (bc *Client) doRequest(req *http.Request, data interface{}) (*FwdData, error) {
	var (
		errStatus ErrBinanceStatus
	)
	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the request: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}
	_ = resp.Body.Close()
	fwd := &FwdData{
		Status:      resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Data:        respBody,
	}
	switch resp.StatusCode {
	case http.StatusOK:
		if data == nil { // if data == nil then caller does not care about response body, consider as success
			return fwd, nil
		}
		if err = json.Unmarshal(respBody, data); err != nil {
			return fwd, fmt.Errorf("failed to parse data into struct: %s", respBody)
		}
	default:
		if err = json.Unmarshal(respBody, &errStatus); err != nil {
			return fwd, fmt.Errorf("failed to parse data into struct: %s", respBody)
		}
		return fwd, fmt.Errorf("binance response with error: %w", &errStatus)
	}
	return fwd, nil
}

func (bc *Client) AllCoinInfo() ([]common.CoinInfo, error) {
	var result []common.CoinInfo

	requestURL := fmt.Sprintf("%s/sapi/v1/capital/config/getall", apiBaseURL)
	req, err := NewRequestBuilder(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	rr := req.WithHeader(apiKeyHeader, bc.apiKey).
		SignedRequest(bc.secretKey)
	_, err = bc.doRequest(rr, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
