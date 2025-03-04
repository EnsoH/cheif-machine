package exchangeAdapters

import (
	"cw/logger"
	"fmt"
	"log"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type MexcAdapter struct {
	Client *ccxt.Mexc
}

func (m *MexcAdapter) Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error) {
	return m.Client.Withdraw(token, amount, address, ccxt.WithWithdrawParams(map[string]interface{}{
		"network": network,
	}))
}

func (m *MexcAdapter) GetBalance(symbol string) (ccxt.Balances, error) {
	return m.Client.FetchBalance()
}

func (m *MexcAdapter) GetPrices(symbol string) (float64, error) {
	m.Client.LoadMarkets()
	ticker, err := m.Client.FetchTicker(fmt.Sprintf("%s/USDT", symbol))
	if err != nil {
		return 0, err
	}
	return *ticker.Last, nil
}

func (m *MexcAdapter) GetChains(token string) error {
	marketsChan := m.Client.LoadMarkets()
	result := <-marketsChan

	if err, ok := result.(error); ok && err != nil {
		logger.GlobalLogger.Error(fmt.Sprintf("failed to load markets: %+v", err))
		return fmt.Errorf("failed to load markets: %w", err)
	}

	curRaw, ok := m.Client.Currencies[token]
	if !ok {
		return fmt.Errorf("token %s not found", token)
	}

	log.Printf("info: %v", curRaw)
	return nil
}
