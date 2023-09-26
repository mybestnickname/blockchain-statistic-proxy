package coingeckometrics

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	coingeckoprovider "external-metrics/pkg/coingecko"
	"external-metrics/pkg/httperror"
	"external-metrics/pkg/tools/logging"

	"github.com/gin-gonic/gin"
)

type GetCoinChartsReq struct {
	CoinShort  string `form:"coin" binding:"required"`
	ConvCurr   string `form:"conversion"`
	RangeStart string `form:"rangeStart"`
	RangeEnd   string `form:"rangeEnd"`
}

// GetCoinChart получение данных для построения графиков цены
// Data granularity is automatic
// 1 day from current time = 5 minute interval data
// 1 - 90 days from current time = hourly data
// above 90 days from current time = daily data (00:00 UTC)
func GetCoinChart(coingeckoProvider *coingeckoprovider.Provider, logger *logging.Logger) func(c *gin.Context) {
	return httperror.ErrorWrapper(logger, func(c *gin.Context) (interface{}, int, error) {
		logger.Infof("Start GetCoinChart...")

		coinChartReq, isErrored := parseGetCoinChartRequest(c, logger)
		if isErrored {
			return nil, http.StatusBadRequest, fmt.Errorf("query args parsing error")
		}

		logger.Infof("Parse request successfully")
		logger.Debug("Request is %+v", coinChartReq)

		coinChartResp, err := coingeckoProvider.GetCoinGeckoCoinChart(
			c,
			coinChartReq.CoinShort,
			coinChartReq.ConvCurr,
			coinChartReq.RangeStart,
			coinChartReq.RangeEnd,
		)
		if err != nil {
			logger.Errorf(
				fmt.Sprintf("Can't get coin chart for %s %s", coinChartReq.CoinShort, coinChartReq.ConvCurr),
			)
			return nil, http.StatusNotFound, fmt.Errorf(
				fmt.Sprintf("failed to get market chart for %s in %s", coinChartReq.CoinShort, coinChartReq.ConvCurr),
			)
		}

		logger.Infof("GetCoinGeckoCoinChart successfully")

		return coinChartResp, http.StatusOK, nil
	})
}

func parseGetCoinChartRequest(
	c *gin.Context,
	logger *logging.Logger,
) (coinChartReq GetCoinChartsReq, isErrored bool) {
	isErrored = true

	if err := c.ShouldBindQuery(&coinChartReq); err != nil {
		logger.Errorf("Missing parameter(s) in url. %v", err)
		return
	}

	if coinChartReq.ConvCurr == "" {
		coinChartReq.ConvCurr = "usd"
	}

	// default range is 1 day
	currentTime := time.Now()
	rangeStart, err := strconv.ParseInt(coinChartReq.RangeStart, 10, 64)
	if err != nil || rangeStart > currentTime.Unix() || rangeStart < 0 {
		coinChartReq.RangeStart = strconv.FormatInt(currentTime.AddDate(0, 0, -1).Unix(), 10)
	}
	rangeEnd, err := strconv.ParseInt(coinChartReq.RangeEnd, 10, 64)
	if err != nil || rangeEnd > currentTime.Unix() || rangeEnd < 0 {
		coinChartReq.RangeEnd = strconv.FormatInt(currentTime.Unix(), 10)
	}

	return coinChartReq, false
}
