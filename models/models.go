package models

type Config struct {
	Threads    int       `json:"threads"`
	CEXConfigs CexConfig `json:"cex"`
}

type CexConfig struct {
	BybitCfg Bybit `json:"bybit"`
	// Cex        string `json:"cex"`
	// API_key    string `json:"api_key"`
	// API_secret string `json:"secret_key"`
	// Chain       []string `json:"chain"`
	// Currency    []string `json:"currency"`
	// AmountRange []int    `json:"amount_range"`
	// Decimals    int      `json:"decimals"`
}

type Bybit struct {
	API_key          string `json:"api_key"`
	API_secret       string `json:"secret_key"`
	WithdrawEndpoint string `json:"withdraw_endpoint"`
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
