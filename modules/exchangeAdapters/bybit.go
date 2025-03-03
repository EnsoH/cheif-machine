package exchangeAdapters

import (
	"fmt"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type BybitAdapter struct {
	Client *ccxt.Bybit
}

func (ba *BybitAdapter) Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error) {
	return ba.Client.Withdraw(token, amount, address, ccxt.WithdrawOptions(ccxt.WithWithdrawParams(map[string]interface{}{
		"forceChain": 1,
		"network":    network,
	})))
}

func (ba *BybitAdapter) GetBalance(symbol string) (ccxt.Balances, error) {
	return ba.Client.FetchBalance(map[string]interface{}{
		"type": "funding",
	})
}

func (ba *BybitAdapter) GetPrices(symbol string) (float64, error) {
	ba.Client.LoadMarkets()
	ticker, err := ba.Client.FetchTicker(fmt.Sprintf("%s/USDT", symbol))
	if err != nil {
		return 0, err
	}
	return *ticker.Last, nil
}
