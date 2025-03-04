package exchangeAdapters

import (
	"fmt"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type BinanceAdapter struct {
	Client *ccxt.Binance
}

func (b *BinanceAdapter) Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error) {
	return b.Client.Withdraw(token, amount, address, ccxt.WithdrawOptions(ccxt.WithWithdrawParams(map[string]interface{}{
		"network": network,
	})))
}

func (b *BinanceAdapter) GetBalance(symbol string) (ccxt.Balances, error) {
	return b.Client.FetchBalance(map[string]interface{}{
		"type": "spot",
	})
}

func (b *BinanceAdapter) GetPrices(symbol string) (float64, error) {
	b.Client.LoadMarkets()
	ticker, err := b.Client.FetchTicker(fmt.Sprintf("%s/USDT", symbol))
	if err != nil {
		return 0, err
	}
	return *ticker.Last, nil
}
