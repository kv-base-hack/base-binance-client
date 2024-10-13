package worker

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	inmem "github.com/kv-base-hack/common/inmem_db"
	"github.com/kv-base-hack/kv-binance/internal/storage"
	"github.com/kv-base-hack/kv-binance/internal/util"
	"go.uber.org/zap"
)

const PricesKey = "binance_prices"
const MarkPricesKey = "binance_mark_prices"

type Handler struct {
	log                       *zap.SugaredLogger
	store                     *storage.Storage
	getPriceDuration          time.Duration
	bSpot                     *binance.Client
	bFuture                   *futures.Client
	updateBinanceInfoDuration time.Duration
	inMemDB                   inmem.Inmem
}

func NewHandler(log *zap.SugaredLogger, store *storage.Storage,
	getPriceDuration time.Duration, bSpot *binance.Client, bFuture *futures.Client,
	updateBinanceInfoDuration time.Duration, inMemDB inmem.Inmem) *Handler {
	return &Handler{
		log:                       log,
		store:                     store,
		getPriceDuration:          getPriceDuration,
		bSpot:                     bSpot,
		bFuture:                   bFuture,
		updateBinanceInfoDuration: updateBinanceInfoDuration,
		inMemDB:                   inMemDB,
	}
}

// func (h *Handler) UpdatePrices() {
// 	log := h.log.With("worker", "get_binance_prices")
// 	ticker := time.NewTicker(h.getPriceDuration)

// 	for ; ; <-ticker.C {
// 		askPrice := make(map[string]string)
// 		bidPrice := make(map[string]string)

// 		prices, err := h.bFuture.NewListBookTickersService().Do(context.Background())
// 		if err != nil {
// 			log.Errorw("error when get binance prices", "err", err)
// 			continue
// 		}
// 		for _, p := range prices {
// 			askPrice[p.Symbol] = p.AskPrice
// 			bidPrice[p.Symbol] = p.BidPrice
// 		}
// 		h.store.SetPrice(askPrice, bidPrice)
// 		pricesBin, _ := json.Marshal(prices)
// 		err = h.inMemDB.Set(PricesKey, pricesBin, time.Second*5)
// 		if err != nil {
// 			log.Errorw("error when set inmem db", "err", err)
// 		}

// 		markPrices := make(map[string]string)
// 		premium := make(map[string]string)
// 		mPrices, err := h.bFuture.NewPremiumIndexService().Do(context.Background())
// 		if err != nil {
// 			log.Error("error when get binance mark prices", "err", err)
// 			continue
// 		}
// 		for _, p := range mPrices {
// 			markPrices[p.Symbol] = p.MarkPrice
// 			premium[p.Symbol] = p.LastFundingRate
// 		}
// 		h.store.SetMarkPrice(markPrices)
// 		mPricesBin, _ := json.Marshal(mPrices)
// 		err = h.inMemDB.Set(MarkPricesKey, mPricesBin, time.Second*5)
// 		if err != nil {
// 			log.Errorw("error when set inmem db", "err", err)
// 		}
// 		h.store.SetPremium(premium)
// 	}
// }

func (h *Handler) UpdateBinanceFuturePairs() {
	log := h.log.With("worker", "get_binance_future_pairs")
	ticker := time.NewTicker(time.Hour * 8)

	for ; ; <-ticker.C {
		symbols := []string{}
		pairs, err := h.bFuture.NewExchangeInfoService().Do(context.Background())
		if err != nil {
			log.Errorw("error when get binance prices", "err", err)
			continue
		}
		for _, p := range pairs.Symbols {
			if !strings.HasSuffix(p.Symbol, "USDT") {
				continue
			}
			symbols = append(symbols, p.Symbol)
		}
		h.store.SetFutureSymbol(symbols)
	}
}

func (h *Handler) UpdateBinanceSpotPairs() {
	log := h.log.With("worker", "get_binance_spot_pairs")
	ticker := time.NewTicker(time.Hour * 8)

	for ; ; <-ticker.C {
		symbols := []string{}
		pairs, err := h.bSpot.NewExchangeInfoService().Do(context.Background())
		if err != nil {
			log.Errorw("error when get binance prices", "err", err)
			continue
		}
		for _, p := range pairs.Symbols {
			if !strings.HasSuffix(p.Symbol, "USDT") {
				continue
			}
			if p.Status != "TRADING" {
				continue
			}
			spotPair := false
			for _, permission := range p.Permissions {
				if permission == "SPOT" {
					spotPair = true
					break
				}
			}
			if !spotPair {
				continue
			}
			symbols = append(symbols, p.Symbol)
		}
		h.store.SetSpotSymbolWithUsdt(symbols)
	}
}

// get value at pos of s, separate by sep. index from 1
// func (h *Handler) getValueFromString(s string, pos int, sep string) (string, bool) {
// 	values := strings.Split(s, sep)
// 	if len(values) < pos {
// 		return "", false
// 	}
// 	return strings.ReplaceAll(values[pos-1], "/", ""), true
// }

// func (h *Handler) handleNewSignal(log *zap.SugaredLogger, signal *db.TelegramSignal, price string, users []db.User) {
// 	for _, u := range users {
// 		ciphertextByte, err := hex.DecodeString(u.BinanceApiSecret)
// 		if err != nil {
// 			log.Errorw("invalid cannot decode string", "u.BinanceApiSecret", u.BinanceApiSecret, "err", err)
// 			continue
// 		}
// 		secretByte, err := h.enc.Decrypt(ciphertextByte)
// 		if err != nil {
// 			log.Errorw("invalid request cannot decrypt", "ciphertextByte", ciphertextByte, "err", err)
// 			continue
// 		}
// 		_, bFuture := bf.NewBinance(u.BinanceApiKey, string(secretByte))
// 		bOrder, err := order.NewBinanceOrder(log, bFuture, int(u.MaxOpenPositions), u.Email, h.store)
// 		if err != nil {
// 			log.Errorw("cannot init binanceOrder", "ciphertextByte", ciphertextByte, "err", err)
// 			continue
// 		}
// 		// add 4 bps
// 		if signal.Action == common.SignalActionLong {
// 			price, _ = util.MultiplyString(price, "1.0004")
// 		} else {
// 			price, _ = util.MultiplyString(price, "0.9996")
// 		}

// 		// 250 bps
// 		stlPrice, _ := util.MultiplyString(price, "0.975")
// 		if signal.Action == common.SignalActionShort {
// 			stlPrice, _ = util.MultiplyString(price, "1.025")
// 		}
// 		bOrder.CreateOrder(signal.Symbol, signal.Action, price, fmt.Sprintf("%.6f", u.UsdtOrderLimit), stlPrice)
// 	}
// 	msg := "👀 kaivest Algo Notification 👀\n"
// 	msg += signal.Symbol + " All entry targets achieved\n"
// 	msg += fmt.Sprintf("Average Entry Price: %s 💵\n", price)
// 	// TODO: check mark prices with entry prices, consider cancel by exit entry zone
// 	err := h.botMsg.SendReplyMsg(signal.TelegramGroupID, msg, signal.SignalID)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "replied message not found") {
// 			signal.Status = common.SignalStatusCancelByRemove
// 		}
// 		log.Errorw("error when send msg to update signal status", "signal.Symbol", signal.Symbol, "signal.SignalID", signal.SignalID, "err", err)
// 	}
// 	signal.Status = common.SignalStatusNotified
// }

// func (h *Handler) handleNotifiedSignal(log *zap.SugaredLogger, signal *db.TelegramSignal, price string, duration time.Duration, users []db.User) {
// 	entry, err := strconv.ParseFloat(signal.Entry, 64)
// 	if err != nil {
// 		log.Errorw("error when parse entry price", "err", err)
// 		return
// 	}
// 	stoploss, err := strconv.ParseFloat(signal.Stoploss, 64)
// 	if err != nil {
// 		log.Errorw("error when parse stop loss", "err", err)
// 		return
// 	}
// 	mPrice, err := strconv.ParseFloat(price, 64)
// 	if err != nil {
// 		log.Errorw("error when parse mark price", "symbol", signal.Symbol, "price", price, "err", err)
// 		return
// 	}

// 	// check profit or stoploss
// 	if (signal.Action == common.SignalActionLong && mPrice < stoploss) || (signal.Action == common.SignalActionShort && mPrice > stoploss) {
// 		for _, u := range users {
// 			ciphertextByte, err := hex.DecodeString(u.BinanceApiSecret)
// 			if err != nil {
// 				log.Errorw("invalid cannot decode string", "u.BinanceApiSecret", u.BinanceApiSecret, "err", err)
// 				continue
// 			}
// 			secretByte, err := h.enc.Decrypt(ciphertextByte)
// 			if err != nil {
// 				log.Errorw("invalid request cannot decrypt", "ciphertextByte", ciphertextByte, "err", err)
// 				continue
// 			}
// 			_, bFuture := bf.NewBinance(u.BinanceApiKey, string(secretByte))
// 			bOrder, err := order.NewBinanceOrder(log, bFuture, int(u.MaxOpenPositions), u.Email, h.store)
// 			if err != nil {
// 				log.Errorw("cannot init binanceOrder", "ciphertextByte", ciphertextByte, "err", err)
// 				continue
// 			}
// 			// close all position
// 			bOrder.ClosePosition(signal.Symbol, price, "1")
// 		}
// 		msg := "👀 kaivest Algo Notification 👀\n"
// 		if signal.LastProfitPrice != "" {
// 			msg += signal.Symbol + " Closed at trailing stoploss after reaching take profit ⚠️\n"
// 			signal.Status = common.SignalStatusCancelByTrailingStoploss
// 		} else {
// 			msg += signal.Symbol + " Stoploss ⛔️\n"
// 			msg += fmt.Sprintf("Loss: %.2f", math.Abs((mPrice-entry)/entry)*100) + "%\n"
// 			signal.Status = common.SignalStatusCancelByStoploss
// 		}
// 		err := h.botMsg.SendReplyMsg(signal.TelegramGroupID, msg, signal.SignalID)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "replied message not found") {
// 				signal.Status = common.SignalStatusCancelByRemove
// 			}
// 			log.Errorw("error when send msg to update signal status", "signal.Symbol", signal.Symbol, "signal.SignalID", signal.SignalID, "err", err)
// 		}
// 		signal.LastStoploss = price
// 		return
// 	}
// 	tpBytes, err := signal.TakeProfits.MarshalJSON()
// 	if err != nil {
// 		log.Errorw("error when marshal take profits", "err", err)
// 		return
// 	}
// 	var takeProfits []common.TakeProfit
// 	if err := json.Unmarshal(tpBytes, &takeProfits); err != nil {
// 		log.Errorw("error when un marshal take profit", "err", err)
// 		return
// 	}
// 	isCompleted := false
// 	for i, tp := range takeProfits {
// 		if tp.Status == common.SignalTakeProfitStatusNotified {
// 			continue
// 		}
// 		tpPrice, err := strconv.ParseFloat(tp.Price, 64)
// 		if err != nil {
// 			log.Errorw("error when parse stop loss", "err", err)
// 			continue
// 		}
// 		// check profit or stoploss
// 		if (signal.Action == common.SignalActionLong && mPrice > tpPrice) || (signal.Action == common.SignalActionShort && mPrice < tpPrice) {
// 			for _, u := range users {
// 				ciphertextByte, err := hex.DecodeString(u.BinanceApiSecret)
// 				if err != nil {
// 					log.Errorw("invalid cannot decode string", "u.BinanceApiSecret", u.BinanceApiSecret, "err", err)
// 					continue
// 				}
// 				secretByte, err := h.enc.Decrypt(ciphertextByte)
// 				if err != nil {
// 					log.Errorw("invalid request cannot decrypt", "ciphertextByte", ciphertextByte, "err", err)
// 					continue
// 				}
// 				_, bFuture := bf.NewBinance(u.BinanceApiKey, string(secretByte))
// 				bOrder, err := order.NewBinanceOrder(log, bFuture, int(u.MaxOpenPositions), u.Email, h.store)
// 				if err != nil {
// 					log.Errorw("cannot init binanceOrder", "ciphertextByte", ciphertextByte, "err", err)
// 					continue
// 				}
// 				// close all position
// 				bOrder.ClosePosition(signal.Symbol, price, "0.2")
// 			}

// 			msg := "👀 kaivest Algo Notification 👀\n"
// 			msg += signal.Symbol + fmt.Sprintf(" Take-Profit %d ✅\n", i+1)
// 			msg += fmt.Sprintf("Profit: %.2f", math.Abs((mPrice-entry)/entry*100)) + "%\n"
// 			msg += fmt.Sprintf("Period: %s  ⏰\n", util.GetDurationFormat(duration))
// 			err := h.botMsg.SendReplyMsg(signal.TelegramGroupID, msg, signal.SignalID)
// 			if err != nil {
// 				if strings.Contains(err.Error(), "replied message not found") {
// 					signal.Status = common.SignalStatusCancelByRemove
// 				}
// 				log.Errorw("error when send msg to update signal status", "signal.Symbol", signal.Symbol, "signal.SignalID", signal.SignalID, "err", err)
// 			}
// 			// move stoploss to entry or last take profit
// 			// stl entry 0 1 2 3 ....
// 			if i == 1 {
// 				signal.Stoploss = signal.Entry
// 			} else if i > 1 {
// 				signal.Stoploss = takeProfits[i-2].Price
// 			}
// 			signal.LastProfitPrice = price
// 			takeProfits[i].Status = common.SignalTakeProfitStatusNotified
// 			takeProfits[i].UpdatedAt = time.Now()
// 			if i == len(takeProfits)-1 {
// 				isCompleted = true
// 			}
// 		}
// 		// TODO: call binance here to cancel all related
// 	}

// 	takeProfitByte, _ := json.Marshal(takeProfits)
// 	signal.TakeProfits = takeProfitByte

// 	if isCompleted {
// 		signal.Status = common.SignalStatusCompleted
// 	}
// }

// func (h *Handler) handleCancelByAdminSignal(log *zap.SugaredLogger, signal *db.TelegramSignal, price string) {
// 	msg := "👀 kaivest Algo Notification 👀\n"
// 	msg += signal.Symbol + " Manually Cancelled by Admin"
// 	err := h.botMsg.SendReplyMsg(signal.TelegramGroupID, msg, signal.SignalID)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "replied message not found") {
// 			signal.Status = common.SignalStatusCancelByRemove
// 		}
// 		log.Errorw("error when send msg to cancel signal by admin", "signal.Symbol", signal.Symbol, "signal.SignalID", signal.SignalID, "err", err)
// 	}
// 	signal.Status = common.SignalStatusCancelByAdminNotified
// }

// func (h *Handler) UpdateSignalAndCreateOrder() {
// 	log := h.log.With("worker", "UpdateSignalAndCreateOrder").With("id", utils.RandomString(29))
// 	newSignals := h.db.GetSignal([]common.SignalStatus{common.SignalStatusNew,
// 		common.SignalStatusNotified,
// 		common.SignalStatusCancelByAdmin})
// 	if len(newSignals) == 0 {
// 		log.Debugw("no new signals from db")
// 		return
// 	}

// 	markPrices := h.store.GetMarkPrices()

// 	users, err := h.db.GetEnableBinanceUsers()
// 	if err != nil {
// 		log.Errorw("invalid when get binance users", "err", err)
// 		return
// 	}

// 	for i, s := range newSignals {
// 		if time.Since(s.CreatedAt) > time.Hour*48 && s.LastProfitPrice == "" && s.LastStoploss == "" {
// 			newSignals[i].Status = common.SignalStatusCancelByExpired
// 			continue
// 		}
// 		// TODO: we use send msg sequentially to send msg and update to database, it slows and bad performance.
// 		// consider other ways when scale
// 		if s.Status == common.SignalStatusNew {
// 			h.handleNewSignal(log, &newSignals[i], markPrices[s.Symbol], users)
// 		} else if s.Status == common.SignalStatusNotified {
// 			h.handleNotifiedSignal(log, &newSignals[i], markPrices[s.Symbol], time.Since(s.CreatedAt), users)
// 		} else if s.Status == common.SignalStatusCancelByAdmin {
// 			h.handleCancelByAdminSignal(log, &newSignals[i], markPrices[s.Symbol])
// 		}
// 	}
// 	err = h.db.UpdateSignals(newSignals)
// 	if err != nil {
// 		log.Errorw("error when update signals status. can be inconsistent here", "err", err)
// 	}
// }

func (h *Handler) UpdateBinanceInfo() {
	log := h.log.With("workerID", "UpdateBinanceInfo")
	ticker := time.NewTicker(h.updateBinanceInfoDuration)
	for ; ; <-ticker.C {
		info, err := h.bFuture.NewExchangeInfoService().Do(context.Background())
		if err != nil {
			continue
		}
		pricePrecision := map[string]int{}
		quantityPrecision := map[string]int{}
		notional := map[string]float64{}

		for _, c := range info.Symbols {
			if c.LotSizeFilter() == nil || c.PriceFilter() == nil {
				log.Errorw("empty data. ignore",
					"LotSizeFilter", c.LotSizeFilter(),
					"PriceFilter", c.PriceFilter())
				continue
			}
			minNotionalStr := &futures.MinNotionalFilter{}
			for _, filter := range c.Filters {
				if filter["filterType"].(string) == string("MIN_NOTIONAL") {
					if i, ok := filter["notional"]; ok {
						minNotionalStr.Notional = i.(string)
					}
				}
			}

			quantityPrecision[c.Symbol] = util.PrecisionFromStepSize(c.LotSizeFilter().StepSize)
			pricePrecision[c.Symbol] = util.PrecisionFromStepSize(c.PriceFilter().TickSize)

			minNotional, err := strconv.ParseFloat(minNotionalStr.Notional, 64)
			if err != nil {
				log.Errorw("error when parse min notional", "minNotional", minNotional, "err", err)
				continue
			}
			notional[c.Symbol] = minNotional
		}
		h.store.SetBoq(pricePrecision, quantityPrecision, notional)
		h.store.SetBinanceInfo(info)
	}
}

func (h *Handler) UpdateBinanceSpotInfo() {
	ticker := time.NewTicker(h.updateBinanceInfoDuration)
	for ; ; <-ticker.C {
		info, err := h.bSpot.NewExchangeInfoService().Do(context.Background())
		if err != nil {
			continue
		}
		h.store.SetBinanceSpotInfo(info)
	}
}

func (h *Handler) Run() {
	go h.UpdateBinanceInfo()
	go h.UpdateBinanceSpotInfo()
	go func() {
		h.UpdateBinanceFuturePairs()
	}()
	go func() {
		h.UpdateBinanceSpotPairs()
	}()
	// h.UpdatePrices()
}
