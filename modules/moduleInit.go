package modules

import (
	"cw/config"
	"cw/httpClient"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

// Interface for Exchange Adapter.
type ExchangeModule interface {
	Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error)
	GetBalance(symbol string) (ccxt.Balances, error)
	GetPrices(token string) (float64, error)
	GetChains(token string) error // TODO: реализовать получение списка сетей.
}

// Factory for future scaling using massive simultaneous output from multiple exchanges for more efficient initialisation.
// type ExchangeFactory func() (ExchangeModule, error)

func ModulesInit(cexNames ...string) (*Exchanges, error) {
	hc, err := httpClient.NewHttpClient(
		httpClient.WithHttp2(),
		httpClient.WithProxy(config.Cfg.IpAddresses[0]),
	)
	if err != nil {
		return nil, err
	}

	modules, err := NewExchangeModule(hc, cexNames...)
	if err != nil {
		return nil, err
	}
	return modules, nil
}
