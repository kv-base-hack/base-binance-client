package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/gin-gonic/gin"
	"github.com/kv-base-hack/common/utils"
)

type SpotKLineRequest struct {
	Pair     string `form:"pair" binding:"required"`
	Interval string `form:"interval" binding:"required"`
	Limit    int    `form:"limit" binding:"required"`
	StartTs  int64  `form:"start_ts"`
}

func (s *Server) getSpotKLine(c *gin.Context) {
	log := s.log.With("ID", utils.RandomString(29))

	var request SpotKLineRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		log.Errorw("invalid request when get spot kline", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidKLineRequest.Error()})
		return
	}
	var err error
	var res []*binance.Kline
	for i := 0; i < retry; i++ {
		kLineService := s.bSpot.NewKlinesService().Symbol(request.Pair).Interval(request.Interval).Limit(request.Limit)
		if request.StartTs > 0 {
			kLineService = kLineService.StartTime(request.StartTs)
		}
		res, err = kLineService.Do(context.Background())
		if err != nil {
			log.Errorw("error when get spot kline", "retry", i, "err", err)
			continue
		}
		break
	}
	if err != nil {
		log.Errorw("couldn't spot kline", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) getAllCoinInfo(c *gin.Context) {
	log := s.log.With("ID", utils.RandomString(29))
	res, err := s.bCustomClient.AllCoinInfo()
	if err != nil {
		log.Errorw("error when get all coin info", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

type AllBookTickerRequest struct {
	Symbols string `form:"symbols" binding:"required"`
}

func (s *Server) getAllBookTicker(c *gin.Context) {
	log := s.log.With("ID", utils.RandomString(29))

	var request AllBookTickerRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		log.Errorw("invalid request when get all book tikcer", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidBookTicker.Error()})
		return
	}
	symbols := strings.Split(request.Symbols, ",")
	res, err := s.bSpot.NewListSymbolTickerService().Symbols(symbols).Do(context.Background())
	if err != nil {
		log.Errorw("error when get all coin info", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) getSpotExchangeInfo(c *gin.Context) {
	log := s.log.With("ID", utils.RandomString(29))
	info, err := s.store.GetBinanceSpotInfo()
	if err != nil {
		log.Errorw("error when get binance info", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error when binance info": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (s *Server) getSpotPairsWithUsdt(c *gin.Context) {
	pairs := s.store.GetSpotSymbolsWithUsdt()
	c.JSON(http.StatusOK, pairs)
}
