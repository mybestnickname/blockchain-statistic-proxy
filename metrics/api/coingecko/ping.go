package coingeckometrics

import (
	"fmt"
	"net/http"

	coingeckoprovider "external-metrics/pkg/coingecko"
	"external-metrics/pkg/httperror"
	"external-metrics/pkg/tools/logging"

	"github.com/gin-gonic/gin"
)

func PingCoinGecko(coingeckoProvider *coingeckoprovider.Provider, logger *logging.Logger) func(c *gin.Context) {
	return httperror.ErrorWrapper(logger, func(c *gin.Context) (interface{}, int, error) {
		logger.Infof("Start PingCoinGecko...")
		pingRes, err := coingeckoProvider.PingCoinGeckoApi(c)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("internal server error")
		}

		return pingRes, http.StatusOK, nil
	})
}
