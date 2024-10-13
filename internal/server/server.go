package server

import (
	"fmt"
	"net/http"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/gin-gonic/gin"
	"github.com/kv-base-hack/kv-binance/internal/bf"
	"github.com/kv-base-hack/kv-binance/internal/storage"
	"go.uber.org/zap"
)

const retry = 5

// Server to serve the service.
type Server struct {
	s        *gin.Engine
	bindAddr string
	log      *zap.SugaredLogger
	//db       db.DB
	store         *storage.Storage
	bSpot         *binance.Client
	bFuture       *futures.Client
	bCustomClient *bf.Client
}

// New returns a new server.
func NewServer(bindAddr string, //db db.DB,
	store *storage.Storage, sClient *binance.Client, bFuture *futures.Client, bCustomClient *bf.Client) *Server {
	engine := gin.New()

	engine.Use(gin.Recovery())

	s := &Server{
		s:        engine,
		log:      zap.S(),
		bindAddr: bindAddr,
		//	db:       db,
		store:         store,
		bSpot:         sClient,
		bFuture:       bFuture,
		bCustomClient: bCustomClient,
	}

	gin.SetMode(gin.DebugMode)
	s.register()

	return s
}

// Run runs server.
func (s *Server) Run() error {
	s.log.Debugw("run in ", "s.bindAddr", s.bindAddr)
	if err := s.s.Run(s.bindAddr); err != nil {
		return fmt.Errorf("run server: %w", err)
	}
	return nil
}

func (s *Server) register() {
	s.s.GET("/debug/pprof/*all", gin.WrapH(http.DefaultServeMux))

	binanceSpot := s.s.Group("/binance/spot")
	binanceSpot.GET("/kline", s.getSpotKLine)
	binanceSpot.GET("/all-coin-info", s.getAllCoinInfo)
	binanceSpot.GET("/book-ticker", s.getAllBookTicker)
	binanceSpot.GET("/spot-exchange-info", s.getSpotExchangeInfo)
	binanceSpot.GET("/spot-pair-with-usdt", s.getSpotPairsWithUsdt)

	binanceFuture := s.s.Group("/binance/future")
	binanceFuture.GET("/premium", s.getPremium)
	binanceFuture.GET("/pairs", s.getBinanceFuturePairs)
	binanceFuture.GET("/kline", s.getKLine)
	binanceFuture.GET("/oi-stats", s.getOIStats)
	binanceFuture.GET("/exchange-info", s.getExchangeInfo)
	binanceFuture.GET("/book-ticker", s.getListBookTicker)
}
