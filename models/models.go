package models

type Config struct {
	Threads    int       `json:"threads"`
	CEXConfigs CexConfig `json:"cex"`
}

type CexConfig struct {
	BybitCfg Bybit `json:"bybit"`
}

type Bybit struct {
	API_key         string `json:"api_key"`
	API_secret      string `json:"secret_key"`
	BalanceEndpoint string `json:"balance_endpoint"`
	TickersEndpoint string `json:"tickers_endpoint"`
}

type WithdrawConfig struct {
	CEX         string    `json:"cex"`
	Chain       []string  `json:"chain"`
	Currency    []string  `json:"currency"`
	AmountRange []float64 `json:"amount_range"`
}

type WithdrawAction struct {
	CEX      string
	Address  string
	Chain    string
	Amount   float64
	Currency string
}

type BybitTickerResponse struct {
	Result struct {
		List []struct {
			LastPrice string `json:"lastPrice"`
		} `json:"list"`
	} `json:"result"`
}

type BybitWithdrawResponse struct {
	Result string `json:"result"`
}
