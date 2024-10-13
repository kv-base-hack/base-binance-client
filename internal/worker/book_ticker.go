package worker

import (
	"time"

	"github.com/kv-base-hack/kv-binance/internal/storage"
	"go.uber.org/zap"
)

const day = 7
const duration = 5 * time.Second

type BookTickerWorker struct {
	log     *zap.SugaredLogger
	symbols []string
	store   *storage.Storage
}

func NewBookTickerWorker(log *zap.SugaredLogger, symbols []string, store *storage.Storage) *BookTickerWorker {
	return &BookTickerWorker{
		log:     log,
		symbols: symbols,
		store:   store,
	}
}

func (c *BookTickerWorker) Run() {
	// get data of 7 days
	startTs := time.Now().Add(-day * time.Hour * 24)
	for _, s := range c.symbols {
		go c.handle(s, startTs)
	}
}

func (c *BookTickerWorker) handle(symbol string, startTs time.Time) {
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		lastTs := c.store.GetLatestKLine(symbol)
		lastTsCloseTime := lastTs.CloseTime
		if lastTsCloseTime < startTs.UnixMilli() {
			lastTsCloseTime = startTs.UnixMilli()
		}

		// kLineService := s.bSpot.NewKlinesService().Symbol(request.Pair).Interval(request.Interval).Limit(request.Limit)
		// if request.StartTs > 0 {
		// 	kLineService = kLineService.StartTime(request.StartTs)
		// }
		// res, err := kLineService.Do(context.Background())
	}
}
