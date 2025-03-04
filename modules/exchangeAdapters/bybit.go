package exchangeAdapters

import (
	"fmt"
	"log"

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

func (b *BybitAdapter) GetChains(token string) error {
	// Загружаем рынки
	if err := b.Client.LoadMarkets(); err != nil {
		return fmt.Errorf("failed to load markets: %v", err)
	}

	// Получаем данные о валюте по токену
	curRaw, ok := b.Client.Currencies[token]
	if !ok {
		return fmt.Errorf("token %s not found", token)
	}

	log.Printf("info: %v", curRaw)
	return nil
}
