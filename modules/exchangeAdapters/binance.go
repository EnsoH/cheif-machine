package exchangeAdapters

import (
	"cw/logger"
	"fmt"
	"log"

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

type ChainInfo struct {
	ChainId        string  // идентификатор сети
	WithdrawEnable bool    // возможность вывода
	WithdrawFee    float64 // комиссия на вывод
	WithdrawMin    float64 // минимальная сумма для вывода
}

func (b *BinanceAdapter) GetChains(token string) error {
	marketsChan := b.Client.LoadMarkets()
	result := <-marketsChan

	if err, ok := result.(error); ok && err != nil {
		logger.GlobalLogger.Error(fmt.Sprintf("failed to load markets: %+v", err))
		return fmt.Errorf("failed to load markets: %w", err)
	}

	curRaw, ok := b.Client.Currencies[token]
	if !ok {
		return fmt.Errorf("token %s not found", token)
	}

	log.Printf("info: %v", curRaw)
	return nil
}
