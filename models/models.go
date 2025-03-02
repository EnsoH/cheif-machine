package models

type Config struct {
	Threads     int       `json:"threads"`
	IpAddresses []string  `json:"ip_addresses"`
	CEXConfigs  CexConfig `json:"cex"`
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
	TimeRange   []float64 `json:"time_range"`
}

type WithdrawAction struct {
	CEX       string
	Address   string
	Chain     string
	Amount    float64
	Currency  string
	TimeRange float64
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

type BybitBalanceResponse struct {
	Result struct {
		Balance []struct {
			Coin          string `json:"coin"`
			Transfer      string `json:"transferBalance"`
			WalletBalance string `json:"walletBalance"`
		} `json:"balance"`
	} `json:"result"`
}
