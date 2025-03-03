package modules

import (
	"cw/httpClient"
	"cw/logger"
	"fmt"
	"log"
	"strconv"
)

type Exchanges struct {
	Tokens     map[string]string
	Decimals   map[string]int
	CEXs       map[string]ExchangeModule
	HttpClient *httpClient.HttpClient
}

func NewExchangeModule(httpclinet *httpClient.HttpClient, cexNames ...string) (*Exchanges, error) {
	ex := &Exchanges{
		HttpClient: httpclinet,
		CEXs:       make(map[string]ExchangeModule),
	}

	for _, cex := range cexNames {
		optionsMap[cex](ex)
	}

	return ex, nil
}

func (e *Exchanges) Withdraw(cexName, token, address, network string, amount float64) error {
	str := strconv.FormatFloat(amount, 'f', 6, 64)
	convertAmount, _ := strconv.ParseFloat(str, 64)

	cex, ok := e.CEXs[cexName]
	if !ok {
		return fmt.Errorf("failed to get cex from adapters")
	}

	tx, err := cex.Withdraw(token, address, network, convertAmount)

	if err != nil {
		log.Println("Ошибка при выводе средств:", err)
		return err
	}

	logger.GlobalLogger.Infof("[%s] Withdraw %s успешен. Chain %s. Amount %f TxId: %v", address, token, network, amount, tx.TxId)
	return nil
}

func (e *Exchanges) GetBalances(cexName, token string) (float64, error) {
	cex, ok := e.CEXs[cexName]
	if !ok {
		return 0, fmt.Errorf("failed to get cex from adapters")
	}

	bal, err := cex.GetBalance(token)
	if err != nil {
		return 0, nil
	}

	usdt, err := bal.GetBalance(token)
	if err != nil {
		return 0, nil
	}

	return *usdt.Free, nil

}

func (e *Exchanges) GetPrices(cexName, token string) (float64, error) {
	cex, ok := e.CEXs[cexName]
	if !ok {
		return 0, fmt.Errorf("failed to get cex from adapters")
	}

	if token == "USDT" || token == "USDC" {
		return 1.0, nil
	}

	price, err := cex.GetPrices(token)
	if err != nil {
		logger.GlobalLogger.Error(err)
		return 0, err
	}

	return price, nil
}

// func (ex *Exchange) getExchange() {

// }
