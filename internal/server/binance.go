package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kv-base-hack/common/utils"
)

func (s *Server) getPremium(c *gin.Context) {
	premium := s.store.GetPremium()
	c.JSON(http.StatusOK, premium)
}

func (s *Server) getBinanceFuturePairs(c *gin.Context) {
	symbols := s.store.GetFutureSymbol()
	c.JSON(http.StatusOK, symbols)
}

type KLineRequest struct {
	Pair     string `form:"pair" binding:"required"`
	Interval string `form:"interval" binding:"required"`
	Limit    int    `form:"limit" binding:"required"`
	EndTs    int64  `form:"end_ts"`
}

func (s *Server) getKLine(c *gin.Context) {
	log := s.log.With("ID", utils.RandomString(29))

	var request KLineRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		log.Errorw("invalid request when get kline", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidKLineRequest.Error()})
		return
	}
	kLineService := s.bFuture.NewKlinesService().Symbol(request.Pair).Interval(request.Interval).Limit(request.Limit)
	if request.EndTs > 0 {
		kLineService = kLineService.EndTime(request.EndTs)
	}
	res, err := kLineService.Do(context.Background())
	if err != nil {
		log.Errorw("error when get kline", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(res) != request.Limit {
		log.Errorw("invalid response length", "request.Limit", request.Limit, "res", res)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidKLineLength.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

type OpenInterestStatisticsRequest struct {
	Pair   string `form:"pair" binding:"required"`
	Period string `form:"period" binding:"required"`
	Limit  int    `form:"limit" binding:"required"`
}

func (s *Server) getOIStats(c *gin.Context) {
	log := s.log.With("ID", utils.RandomString(29))

	var request OpenInterestStatisticsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		log.Errorw("invalid request when get OI stats", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidOIStatsRequest.Error()})
		return
	}
	oiStats := s.bFuture.NewOpenInterestStatisticsService().Symbol(request.Pair).Period(request.Period).Limit(request.Limit)
	res, err := oiStats.Do(context.Background())
	if err != nil {
		log.Errorw("error when get oi stats", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) getExchangeInfo(c *gin.Context) {
	log := s.log.With("ID", utils.RandomString(29))
	info, err := s.store.GetBinanceInfo()
	if err != nil {
		log.Errorw("error when get binance info", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error when binance info": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

type ListBookTickerRequest struct {
	Symbol string `form:"symbol"`
}

func (s *Server) getListBookTicker(c *gin.Context) {
	log := s.log.With("ID", utils.RandomString(29))
	var request ListBookTickerRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		log.Errorw("invalid request when list book ticker", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidOIStatsRequest.Error()})
		return
	}

	bookTicker := s.bFuture.NewListBookTickersService()
	if len(request.Symbol) > 0 {
		bookTicker = bookTicker.Symbol(request.Symbol)
	}
	res, err := bookTicker.Do(context.Background())
	if err != nil {
		log.Errorw("error when get list book ticker", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
