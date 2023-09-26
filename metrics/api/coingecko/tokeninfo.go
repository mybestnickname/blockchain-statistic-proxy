package coingeckometrics

import (
	"fmt"
	"net/http"

	"external-metrics/metrics/models"
	coingeckoprovider "external-metrics/pkg/coingecko"
	"external-metrics/pkg/httperror"
	"external-metrics/pkg/tools/logging"

	"github.com/gin-gonic/gin"
)

type GetDefiTokenInfoReq struct {
	ConvCurr        string `form:"conversion"`
	Network         string `form:"network" binding:"required"`
	ContractAddress string `form:"address" binding:"required"`
}

// GetDefiTokenInfo получение информации о токене по его адресу
func GetDefiTokenInfo(coingeckoProvider *coingeckoprovider.Provider, logger *logging.Logger) func(c *gin.Context) {
	return httperror.ErrorWrapper(logger, func(c *gin.Context) (interface{}, int, error) {
		logger.Infof("Start GetDefiTokenInfo...")

		defiTokenInfoReq, isErrored := parseGetDefiTokenInfoRequest(c, logger)
		if isErrored {
			return nil, http.StatusBadRequest, fmt.Errorf("query args parsing error")
		}

		logger.Infof("Parse request successfully")
		logger.Debug(
			"Request is network: %s, contractAddress: %s, convCurr: %s",
			defiTokenInfoReq.Network,
			defiTokenInfoReq.ContractAddress,
			defiTokenInfoReq.ConvCurr,
		)

		cgDefiCoinInfo, err := coingeckoProvider.GetDefiTokenInfo(
			c,
			defiTokenInfoReq.Network,
			defiTokenInfoReq.ContractAddress,
		)
		if err != nil {
			logger.Errorf(
				fmt.Sprintf("Can't get coingecko token info for %s %s",
					defiTokenInfoReq.Network,
					defiTokenInfoReq.ContractAddress,
				),
			)
			return nil, http.StatusNotFound, fmt.Errorf(
				fmt.Sprintf(
					"failed to get token info for %s %s in %s",
					defiTokenInfoReq.Network,
					defiTokenInfoReq.ContractAddress,
					defiTokenInfoReq.ConvCurr,
				),
			)
		}

		defiCoinInfoResp := &models.CoinGeckoDefiTokenInfoResp{
			Name:   cgDefiCoinInfo.Name,
			Symbol: cgDefiCoinInfo.Symbol,
			Icons:  models.CoinGeckoIconsResp(cgDefiCoinInfo.Icons),
		}

		if price, found := cgDefiCoinInfo.MarketData.CurrentPrice[defiTokenInfoReq.ConvCurr]; found {
			defiCoinInfoResp.Price = price
		} else {
			logger.Infof(
				fmt.Sprintf(
					"Can't get price for netw: %s, addr: %s, conv_cur: %s",
					defiTokenInfoReq.Network,
					defiTokenInfoReq.ContractAddress,
					defiTokenInfoReq.ConvCurr,
				),
			)
		}

		logger.Infof("GetDefiTokenInfo successfully")

		return defiCoinInfoResp, http.StatusOK, nil
	})
}

func parseGetDefiTokenInfoRequest(
	c *gin.Context,
	logger *logging.Logger,
) (defiTokenInfoReq GetDefiTokenInfoReq, isErrored bool) {
	isErrored = true

	if err := c.ShouldBindQuery(&defiTokenInfoReq); err != nil {
		logger.Errorf("Missing parameter(s) in url. %v", err)
		return
	}
	if defiTokenInfoReq.ConvCurr == "" {
		defiTokenInfoReq.ConvCurr = "usd"
	}
	return defiTokenInfoReq, false
}
