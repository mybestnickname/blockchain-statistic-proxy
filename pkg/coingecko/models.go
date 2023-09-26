package coingeckoprovider

type CoinGeckoMarketCapPercentage struct {
	MarketCapPercentage map[string]float32 `json:"market_cap_percentage"`
}

type CoinGeckoGlobalCryptoData struct {
	CoinGeckoMarketCapPercentage `json:"data"`
}

type CoinGeckoCoinData struct {
	Description map[string]string `json:"description" jsonschema:"required"`
}

type CoinGeckoIcons struct {
	Thumb string `json:"thumb"`
	Small string `json:"small"`
	Large string `json:"large"`
}

type CoinGeckoDefiCoinInfo struct {
	Name       string         `json:"name"`
	Symbol     string         `json:"symbol"`
	Icons      CoinGeckoIcons `json:"image"`
	MarketData struct {
		CurrentPrice map[string]float32 `json:"current_price"`
	} `json:"market_data"`
}
