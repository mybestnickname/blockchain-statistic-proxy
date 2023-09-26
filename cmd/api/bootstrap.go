package main

import (
	"external-metrics/config"
	coingeckometrics "external-metrics/metrics/api/coingecko"
	coingeckoprovider "external-metrics/pkg/coingecko"
	"external-metrics/pkg/tools/logging"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginlogrus "github.com/toorop/gin-logrus"
)

func bootstrapAPI(
	coingeckoProvider *coingeckoprovider.Provider,
	cfg *config.Config,
	logger *logging.Logger,
) *http.Server {
	router := gin.New()
	router.Use(ginlogrus.Logger(logger.Logger))

	apiV1 := router.Group("/api/v1")

	attachRoutesAPI(coingeckoProvider, logger, apiV1)

	return &http.Server{
		Handler:      router,
		Addr:         cfg.App.APIAddress,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
}

func attachRoutesAPI(
	coingeckoProvider *coingeckoprovider.Provider,
	logger *logging.Logger,
	router *gin.RouterGroup,
) {
	router.GET("/ping_gecko", coingeckometrics.PingCoinGecko(coingeckoProvider, logger))

	router.GET("/coin/info", coingeckometrics.GetCoinInfo(coingeckoProvider, logger))
	router.GET("/coin/chart", coingeckometrics.GetCoinChart(coingeckoProvider, logger))

	router.GET("/token/info", coingeckometrics.GetDefiTokenInfo(coingeckoProvider, logger))
}

func startServer(srv *http.Server) <-chan error {
	errChan := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		errChan <- err
	}()

	return errChan
}
