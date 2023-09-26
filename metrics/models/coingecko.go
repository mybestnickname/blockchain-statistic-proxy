package models

type CoinGeckoPingResp struct {
	GeckoSays string `json:"gecko_says"`
}

// /coins/markets
type CoinInfoResp struct {
	Name                string            `json:"name"`                        // Имя
	Image               string            `json:"image"`                       // Изображение
	CurrentPrice        float32           `json:"current_price"`               // Курс монеты в валюте conversion
	PriceChange24       float32           `json:"price_change_24h"`            // Изменение за 24 часа в валюте conversion
	PriceChange24Perc   float32           `json:"price_change_percentage_24h"` // Изменение за 24 часа в процентах
	MarketCap           float32           `json:"market_cap"`                  //Капитализация в валюте conversion
	CirculatingSupply   float32           `json:"circulating_supply"`          // Монет в обороте
	TotalSupply         float32           `json:"total_supply"`                // Монет всего
	MaxSupply           float32           `json:"max_supply"`                  // Монет всего
	Rank                int               `json:"market_cap_rank"`             // Ранг монеты
	Ath                 float32           `json:"ath"`                         // Максимум за все время
	Volume24            float32           `json:"volume_24h"`                  // Объем торгов за 24 часа в валюте conversion
	MarketCapPercentage float32           `json:"market_cap_percentage"`       // Доля рынка в процентах
	About               map[string]string `json:"about"`                       // Описание монеты на англ
}

type CoinGeckoCoinChartResp struct {
	Prices [][]float64 `json:"prices"`
}

type CoinGeckoIconsResp struct {
	Thumb string `json:"thumb"`
	Small string `json:"small"`
	Large string `json:"large"`
}

type CoinGeckoDefiTokenInfoResp struct {
	Name   string             `json:"name"`
	Symbol string             `json:"symbol"`
	Icons  CoinGeckoIconsResp `json:"image"` // Ссылки на изображения
	Price  float32            `json:"current_price"`
}
