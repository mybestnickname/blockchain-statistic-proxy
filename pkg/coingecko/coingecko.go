package coingeckoprovider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"external-metrics/metrics/models"
	"external-metrics/pkg/tools/logging"

	cache "github.com/patrickmn/go-cache"
	coingecko "github.com/superoo7/go-gecko/v3"
)

// Provider for coingecko
type Provider struct {
	apiAddress string
	logger     *logging.Logger
	memCache   *cache.Cache
}

// coin struct for cache store
type geckoCoin struct {
	coinId   string
	coinName string
}

// NewProvider return new Coingecko provider object.
func NewProvider(apiAddress string, logger *logging.Logger) (*Provider, error) {
	return &Provider{
		apiAddress: apiAddress,
		logger:     logger,
		memCache:   cache.New(cache.NoExpiration, 10*time.Minute),
	}, nil
}

func (p *Provider) PingCoinGeckoApi(ctx context.Context) (*models.CoinGeckoPingResp, error) {
	p.logger.Infof("Start PingCoinGeckoApi provider method...")

	respBody, err := p.Do(ctx, "GET", "/ping", nil)
	if err != nil {
		p.logger.Errorf("Can't ping coingecko at %s: %v", p.apiAddress, err)
		return nil, err
	}

	p.logger.Infof("Get success response")

	res := &models.CoinGeckoPingResp{}
	if err = json.Unmarshal(respBody, res); err != nil {
		p.logger.Errorf("ReqToCoinGeckoApi: Can not unmarshal coingecko response ping body %s: %v", p.apiAddress, err)
		return res, err
	}

	return res, nil
}

func (p *Provider) GetDefiTokenInfo(
	ctx context.Context,
	network string,
	contractAddress string,
) (*CoinGeckoDefiCoinInfo, error) {
	p.logger.Infof("Start GetDefiTokenInfo provider method...")

	requestURL := fmt.Sprintf("/coins/%s/contract/%s", network, contractAddress)
	respBody, err := p.Do(ctx, "GET", requestURL, nil)
	if err != nil {
		p.logger.Errorf("Can't GET coin info by contract address from coingecko. URL: %s: %v", requestURL, err)

		return nil, err
	}

	p.logger.Infof("Get success response")

	var coinDefiInfo CoinGeckoDefiCoinInfo
	if err = json.Unmarshal(respBody, &coinDefiInfo); err != nil {
		p.logger.Errorf("Can't unmarshal coingecko defiCoinInfo body. URL: %s: %v", requestURL, err)
		return nil, err
	}
	return &coinDefiInfo, nil
}

func (p *Provider) GetCoinGeckoCoinChart(
	ctx context.Context,
	coinShort string,
	convCurr string,
	rangeStart string,
	rangeEnd string,
) (*models.CoinGeckoCoinChartResp, error) {
	p.logger.Infof("Start GetCoinChart provider method...")

	coinId, _, err := p.GetCoinIDName(ctx, coinShort)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"vs_currency": {convCurr},
		"from":        {rangeStart},
		"to":          {rangeEnd},
	}
	requestURL := fmt.Sprintf("/coins/%s/market_chart/range?", coinId) + params.Encode()
	respBody, err := p.Do(ctx, "GET", requestURL, nil)
	if err != nil {
		p.logger.Errorf("Can't GET coin chart from coingecko. URL: %s: %v", requestURL, err)
		return nil, err
	}

	p.logger.Infof("Get success response")

	var coinChart models.CoinGeckoCoinChartResp
	if err = json.Unmarshal(respBody, &coinChart); err != nil {
		p.logger.Errorf("Can't unmarshal coingecko getCoinChart body. URL: %s: %v", requestURL, err)
		return nil, err
	}
	return &coinChart, nil
}

func (p *Provider) GetCoinInfo(
	ctx context.Context,
	coinShort string,
	convCurr string,
) (*models.CoinInfoResp, error) {
	p.logger.Infof("Starting GetCoinInfo provider method...")

	coinId, coinName, err := p.GetCoinIDName(ctx, coinShort)
	if err != nil {
		return nil, err
	}

	// getting coin market info
	params := url.Values{
		"vs_currency": {convCurr},
		"ids":         {coinId},
		"order":       {"market_cap_desc"},
		"per_page":    {"100"},
		"page":        {"1"},
		"sparkline":   {"false"},
	}
	url := "/coins/markets?" + params.Encode()
	respBody, err := p.Do(ctx, "GET", url, nil)
	if err != nil {
		p.logger.Errorf("Can't get coingecko response GET coins market. URL: %s: %v", url, err)
		return &models.CoinInfoResp{}, err
	}
	var coinsInfo []models.CoinInfoResp
	if err = json.Unmarshal(respBody, &coinsInfo); err != nil {
		p.logger.Errorf("Can't unmarshal coingecko GET coins market resp. URL: %s: %v", url, err)
		return &models.CoinInfoResp{}, err
	}

	// set coin name
	coinsInfo[0].Name = coinName

	// FIXME по-хорошему это все сделать в горутины, чтобы параллельно посылать запросы. Но это отдельное веселье, так что по желанию

	// getting last 24h market volume
	volume24h, err := p.Get24hVolume(ctx, coinId, convCurr)
	if err != nil {
		return &coinsInfo[0], nil
	}
	coinsInfo[0].Volume24 = volume24h

	// getting current time market cap percentage
	coinMarketPercentage, err := p.GetMarketCapPercentage(ctx, coinShort)
	if err != nil {
		return &coinsInfo[0], err
	}
	if coinMarketPercentage != 0 {
		coinsInfo[0].MarketCapPercentage = coinMarketPercentage
	}

	// getting coin description(eng)
	coinDescription, err := p.GetCoinDescription(ctx, coinId)
	if err != nil {
		return &coinsInfo[0], err
	}
	coinsInfo[0].About = coinDescription

	return &coinsInfo[0], nil
}

func (p *Provider) GetCoinDescription(ctx context.Context, coinId string) (map[string]string, error) {
	params := url.Values{
		"localization":   {"false"},
		"tickers":        {"false"},
		"market_data":    {"false"},
		"community_data": {"false"},
		"developer_data": {"false"},
		"sparkline":      {"false"},
	}
	url := fmt.Sprintf("/coins/%s?", coinId) + params.Encode()
	respBody, err := p.Do(ctx, "GET", url, nil)
	if err != nil {
		p.logger.Errorf("Can't GET coin info from coingecko. URL: %s: %v", url, err)
		return make(map[string]string), err
	}

	p.logger.Infof("Get success response")

	var description CoinGeckoCoinData
	if err = json.Unmarshal(respBody, &description); err != nil {
		p.logger.Errorf("Can't unmarshal coingecko response GET body. %s %v", url, err)
		return make(map[string]string), err
	}
	return description.Description, nil
}

func (p *Provider) GetMarketCapPercentage(ctx context.Context, coinShort string) (float32, error) {
	respBody, err := p.Do(ctx, "GET", "/global", nil)
	if err != nil {
		p.logger.Errorf("Can't GET global data from coingecko. %v", err)
		return 0, err
	}

	p.logger.Infof("Get success response")
	var globalData CoinGeckoGlobalCryptoData
	err = json.Unmarshal(respBody, &globalData)
	if err != nil {
		p.logger.Errorf("Can't unmarshal coingecko response GET global data body. %v", err)
		return 0, err
	}

	if res, found := globalData.MarketCapPercentage[coinShort]; found {
		return res, nil
	} else {
		p.logger.Info("Market cap percentage in missing in global data for &s", coinShort)
		return 0, nil
	}
}

func (p *Provider) Get24hVolume(ctx context.Context, coinId string, convCurr string) (float32, error) {
	params := url.Values{
		"ids":              {coinId},
		"vs_currencies":    {convCurr},
		"include_24hr_vol": {"true"},
	}
	url := "/simple/price?" + params.Encode()
	respBody, err := p.Do(ctx, "GET", url, nil)
	if err != nil {
		p.logger.Errorf("Can't GET simple price from coingecko. URL: %s: %v", url, err)
		return 0, err
	}

	p.logger.Infof("Get success response")

	var data map[string]map[string]float32
	if err := json.Unmarshal(respBody, &data); err != nil {
		p.logger.Errorf(
			"Can't unmarshal coingecko response GET simple price body for coinId: %s: %v ", coinId, err,
		)
		return 0, err
	}
	h_volume_24 := data[coinId][fmt.Sprintf("%s_24h_vol", convCurr)]
	return h_volume_24, nil
}

func (p *Provider) GetCoinIDName(ctx context.Context, coinShort string) (string, string, error) {
	p.logger.Infof("Starting GetCoinIDName provider method...")

	p.logger.Infof("Trying to get coinId and coinName from cache")
	cacheKey := fmt.Sprintf("%s_coin", coinShort)
	if res, found := p.memCache.Get(cacheKey); found {
		coin := res.(*geckoCoin)
		return coin.coinId, coin.coinName, nil
	}

	p.logger.Infof("Getting from coingecko /coins/list endpoint")
	coinGeckoClient := p.CreateCoinGeckoHttpClient()
	coins, err := coinGeckoClient.CoinsList()
	if err != nil {
		p.logger.Errorf("Can't GET coins list from coingecko for %s: %v", coinShort, err)
		return "", "", err
	}

	for _, coin := range *coins {
		if coinShort == coin.Symbol || coinShort == coin.Name {
			p.logger.Infof("Adding coin id:%s, name:%s into cache.", coin.ID, coin.Name)
			p.memCache.Set(cacheKey, &geckoCoin{coinId: coin.ID, coinName: coin.Name}, cache.DefaultExpiration)

			return coin.ID, coin.Name, nil
		}
	}

	p.logger.Errorf("Can't find coingecko coinId/Name for %s", coinShort)
	return "", "", fmt.Errorf("GetCoinIDName error")
}

func (p *Provider) CreateCoinGeckoHttpClient() *coingecko.Client {
	httpClient := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	coinGeckoClient := coingecko.NewClient(httpClient)
	return coinGeckoClient
}

func (p *Provider) Do(
	ctx context.Context,
	reqType string,
	reqURL string,
	reqBody io.Reader,
) ([]byte, error) {
	httpClient := http.Client{
		Timeout: time.Duration(1) * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, reqType, fmt.Sprintf("%s%s", p.apiAddress, reqURL), reqBody)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Accept", `application/json`)

	resp, err := httpClient.Do(req)
	if err != nil {
		p.logger.Errorf("can't send coingecko request: %v", err)
		return []byte{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.logger.Errorf("can't read coingecko response: %v", err)
		return []byte{}, err
	}

	if resp.StatusCode != 200 {
		p.logger.Errorf("Wrong status code! %d", resp.StatusCode)
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "\t")
		if err != nil {
			p.logger.Errorf("JSON parse error: %v", err)
			return []byte{}, err
		}
		p.logger.Errorf("Error json res: %s", prettyJSON.String())
		return []byte{}, fmt.Errorf("wrong status code resp")
	}

	return body, nil
}
