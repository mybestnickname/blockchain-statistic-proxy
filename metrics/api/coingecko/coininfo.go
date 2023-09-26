package coingeckometrics

import (
	"fmt"
	"net/http"

	coingeckoprovider "external-metrics/pkg/coingecko"
	"external-metrics/pkg/httperror"
	"external-metrics/pkg/tools/logging"

	"github.com/gin-gonic/gin"
)

type GetCoinInfoReq struct {
	CoinShort string `form:"coin" binding:"required"`
	ConvCurr  string `form:"conversion"`
}

// GetCoinInfo получение информации о токене по его имени
func GetCoinInfo(coingeckoProvider *coingeckoprovider.Provider, logger *logging.Logger) func(c *gin.Context) {
	return httperror.ErrorWrapper(logger, func(c *gin.Context) (interface{}, int, error) {
		logger.Infof("Start GetCoinInfo...")

		coinInfoReq, isErrored := parseGetCoinInfoRequest(c, logger)
		if isErrored {
			return nil, http.StatusBadRequest, fmt.Errorf("query args parsing error")
		}

		logger.Infof("Parse request successfully")
		logger.Debug("Request is %+v", coinInfoReq)

		coinInfo, err := coingeckoProvider.GetCoinInfo(c, coinInfoReq.CoinShort, coinInfoReq.ConvCurr)
		if err != nil {
			logger.Errorf(
				fmt.Sprintf("Can't get coingecko coinInfo for %s %s", coinInfoReq.CoinShort, coinInfoReq.ConvCurr),
			)
			return nil, http.StatusNotFound, fmt.Errorf(
				fmt.Sprintf("failed to get coin info for %s with %s conversion", coinInfoReq.CoinShort, coinInfoReq.ConvCurr),
			)
		}

		logger.Infof("GetCoinInfo successfully")

		return coinInfo, http.StatusOK, nil
	})
}

func parseGetCoinInfoRequest(
	c *gin.Context,
	logger *logging.Logger,
) (coinInfoReq GetCoinInfoReq, isErrored bool) {
	isErrored = true

	if err := c.ShouldBindQuery(&coinInfoReq); err != nil {
		logger.Errorf("Missing parameter(s) in url. %v", err)
		return
	}
	if coinInfoReq.ConvCurr == "" {
		coinInfoReq.ConvCurr = "usd"
	}
	return coinInfoReq, false
}
