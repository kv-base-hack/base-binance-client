package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/joho/godotenv"
	inmem "github.com/kv-base-hack/common/inmem_db"
	"github.com/kv-base-hack/common/logger"
	"github.com/kv-base-hack/kv-binance/common"
	"github.com/kv-base-hack/kv-binance/internal/bf"
	"github.com/kv-base-hack/kv-binance/internal/httputil"
	"github.com/kv-base-hack/kv-binance/internal/server"
	"github.com/kv-base-hack/kv-binance/internal/storage"
	"github.com/kv-base-hack/kv-binance/internal/worker"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

const defaultClientTimeout = time.Second * 30

func main() {
	_ = godotenv.Load()
	app := cli.NewApp()
	app.Action = run
	app.Flags = newAppFlags()
	app.Flags = append(app.Flags, logger.NewSentryFlags()...)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.Port)...)
	app.Flags = append(app.Flags, NewPostgreSQLFlags()...)
	app.Flags = append(app.Flags, NewRedisFlags()...)
	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(c *cli.Context) error {
	logger, flusher, err := logger.NewLogger(c)
	if err != nil {
		return err
	}
	defer flusher()

	zap.ReplaceGlobals(logger)
	log := logger.Sugar()
	log.Debugw("Starting application...")
	// database, err := NewDBFromContext(c)
	// if err != nil {
	// 	log.Errorw("error when connect to database", "err", err)
	// 	return err
	// }
	// ps := db.NewPostgres(database)
	store := storage.NewStorage()
	var binanceApiAccounts map[string]common.BinanceAccount
	if err := json.Unmarshal([]byte(c.String(binanceApiKes)), &binanceApiAccounts); err != nil {
		return err
	}
	main := binanceApiAccounts["main"]
	spotClient, futureClient := bf.NewBinance(main.Key, main.Secret)

	redisHost := c.String(redisHostFlag)
	redisPort := c.String(redisPortFlag)
	redisPassword := c.String(redisPasswordFlag)
	redisDB := c.Int(redisDBFlag)
	redisAddr := redisHost + ":" + redisPort
	redis := inmem.NewRedisClient(redisAddr, redisPassword, redisDB)
	handler := worker.NewHandler(log, store,
		c.Duration(getPriceDurationFlag), spotClient, futureClient, c.Duration(updateBinanceInfoDurationFlag), redis)
	go handler.Run()

	httpClient := &http.Client{
		Timeout: defaultClientTimeout,
		Transport: &http.Transport{
			IdleConnTimeout:       time.Second * 120,
			ResponseHeaderTimeout: time.Second * 10,
		},
	}

	bCustomClient := bf.NewClient(main.Key, main.Secret, httpClient)
	host := httputil.NewHTTPAddressFromContext(c)
	server := server.NewServer(host, store, spotClient, futureClient, bCustomClient)
	return server.Run()
}
