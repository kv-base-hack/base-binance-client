package storage

import (
	"sync"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/jinzhu/copier"
)

const defaultPricePrecision = 6
const defaultQuantityPrecision = 6

type BFuture struct {
	// askPrice          map[string]string
	// bidPrice          map[string]string
	// markPrice         map[string]string
	premium           map[string]string
	pricePrecision    map[string]int
	quantityPrecision map[string]int
	notional          map[string]float64
	futureSymbol      []string
	binanceInfo       *futures.ExchangeInfo
}

type BSpot struct {
	symbolsWithUsdt []string
	binanceInfo     *binance.ExchangeInfo
	kLine           map[string][]*binance.Kline
}

type Storage struct {
	mutex   sync.RWMutex
	bFuture BFuture
	bSpot   BSpot
}

func NewStorage() *Storage {
	return &Storage{
		bFuture: BFuture{
			// askPrice:          make(map[string]string),
			// bidPrice:          make(map[string]string),
			// markPrice:         make(map[string]string),
			pricePrecision:    make(map[string]int),
			quantityPrecision: make(map[string]int),
			notional:          make(map[string]float64),
			futureSymbol:      make([]string, 0),
		},
		bSpot: BSpot{
			symbolsWithUsdt: make([]string, 0),
			kLine:           make(map[string][]*binance.Kline),
		},
	}
}

// func (s *Storage) SetPrice(askPrices map[string]string, bidPrices map[string]string) {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()
// 	s.bFuture.askPrice = askPrices
// 	s.bFuture.bidPrice = bidPrices
// }

// func (s *Storage) GetAskPrice(symbol string) string {
// 	s.mutex.RLock()
// 	defer s.mutex.RUnlock()
// 	return s.bFuture.askPrice[symbol]
// }

// func (s *Storage) GetBidPrice(symbol string) string {
// 	s.mutex.RLock()
// 	defer s.mutex.RUnlock()
// 	return s.bFuture.bidPrice[symbol]
// }

// func (s *Storage) GetMapAskPrices() (map[string]string, error) {
// 	s.mutex.RLock()
// 	defer s.mutex.RUnlock()
// 	prices := map[string]string{}
// 	err := copier.CopyWithOption(&prices, s.bFuture.askPrice, copier.Option{DeepCopy: true})
// 	if err != nil {
// 		return prices, err
// 	}
// 	return prices, nil
// }

// func (s *Storage) GetMapBidPrices() (map[string]string, error) {
// 	s.mutex.RLock()
// 	defer s.mutex.RUnlock()
// 	prices := map[string]string{}
// 	err := copier.CopyWithOption(&prices, s.bFuture.bidPrice, copier.Option{DeepCopy: true})
// 	if err != nil {
// 		return prices, err
// 	}
// 	return prices, nil
// }

// func (s *Storage) SetMarkPrice(markPrice map[string]string) {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()
// 	s.bFuture.markPrice = markPrice
// }

func (s *Storage) SetPremium(premium map[string]string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.bFuture.premium = premium
}

// func (s *Storage) GetMarkPrices() map[string]string {
// 	s.mutex.RLock()
// 	defer s.mutex.RUnlock()
// 	prices := map[string]string{}
// 	for k, v := range s.bFuture.markPrice {
// 		prices[k] = v
// 	}
// 	return prices
// }

func (s *Storage) GetPremium() map[string]string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	premium := map[string]string{}
	for k, v := range s.bFuture.premium {
		premium[k] = v
	}
	return premium
}

func (s *Storage) SetBoq(pricePrecision map[string]int, quantityPrecision map[string]int, notional map[string]float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k, v := range pricePrecision {
		s.bFuture.pricePrecision[k] = v
	}
	for k, v := range quantityPrecision {
		s.bFuture.quantityPrecision[k] = v
	}
	for k, v := range notional {
		s.bFuture.notional[k] = v
	}
}

func (s *Storage) GetBinancePricePrecision(symbol string) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	pricePrecision, exist := s.bFuture.pricePrecision[symbol]
	if !exist {
		return defaultPricePrecision
	}
	return pricePrecision
}

func (s *Storage) GetBinanceQuantityPrecision(symbol string) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	quantityPrecision, exist := s.bFuture.quantityPrecision[symbol]
	if !exist {
		return defaultQuantityPrecision
	}
	return quantityPrecision
}

func (s *Storage) GetBinanceNotional() (map[string]float64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	notional := map[string]float64{}
	err := copier.CopyWithOption(&notional, s.bFuture.notional, copier.Option{DeepCopy: true})
	if err != nil {
		return notional, err
	}
	return notional, nil
}

func (s *Storage) SetFutureSymbol(symbols []string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.bFuture.futureSymbol = symbols
}

func (s *Storage) GetFutureSymbol() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.bFuture.futureSymbol
}

func (s *Storage) SetBinanceInfo(info *futures.ExchangeInfo) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.bFuture.binanceInfo = info
}

func (s *Storage) GetBinanceInfo() (futures.ExchangeInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	info := futures.ExchangeInfo{}
	err := copier.CopyWithOption(&info, s.bFuture.binanceInfo, copier.Option{DeepCopy: true})
	if err != nil {
		return info, err
	}
	return info, nil
}

func (s *Storage) SetSpotSymbolWithUsdt(symbols []string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.bSpot.symbolsWithUsdt = symbols
}

func (s *Storage) GetSpotSymbolsWithUsdt() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.bSpot.symbolsWithUsdt
}

func (s *Storage) SetBinanceSpotInfo(info *binance.ExchangeInfo) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.bSpot.binanceInfo = info
}

func (s *Storage) GetBinanceSpotInfo() (binance.ExchangeInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	info := binance.ExchangeInfo{}
	err := copier.CopyWithOption(&info, s.bSpot.binanceInfo, copier.Option{DeepCopy: true})
	if err != nil {
		return info, err
	}
	return info, nil
}

func (s *Storage) AddKline(symbol string, kline []*binance.Kline) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.bSpot.kLine[symbol] = append(s.bSpot.kLine[symbol], kline...)
}

func (s *Storage) GetLatestKLine(symbol string) *binance.Kline {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(s.bSpot.kLine[symbol]) > 0 {
		return s.bSpot.kLine[symbol][len(s.bSpot.kLine)-1]
	}
	return &binance.Kline{}
}

func (s *Storage) GetKline(symbol string, startTs, endTs int64) []*binance.Kline {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	res := []*binance.Kline{}
	for _, t := range s.bSpot.kLine[symbol] {
		if startTs <= t.OpenTime && t.CloseTime <= endTs {
			res = append(res, t)
		}
	}
	return res
}
