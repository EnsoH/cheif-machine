package modules

import (
	"cw/httpClient"
	"cw/logger"
	"cw/models"
	"fmt"
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
		if _, ok := exhangeOptionsMap[cex]; !ok {
			continue
		}

		exhangeOptionsMap[cex](ex)
	}

	return ex, nil
}

func (e *Exchanges) Withdraw(cexName, token, address, network string, amount float64) error {
	// TODO: написать метод для формирования числа с заданным децималс под определенный токен.
	str := strconv.FormatFloat(amount, 'f', 6, 64)
	convertAmount, _ := strconv.ParseFloat(str, 64)

	cex, err := e.getCEX(cexName)
	if err != nil {
		return err
	}

	tx, err := cex.Withdraw(token, address, network, convertAmount)
	if err != nil {
		return err
	}

	logger.GlobalLogger.Infof("[%s] Withdraw %s успешен. Chain %s. Amount %f TxId: %v", address, token, network, amount, tx.TxId)
	return nil
}

func (e *Exchanges) GetBalances(cexName, token string) (float64, error) {
	cex, err := e.getCEX(cexName)
	if err != nil {
		return 0, nil
	}

	bal, err := cex.GetBalance(token)
	if err != nil {
		return 0, nil
	}

	tokenEnable, err := bal.GetBalance(token)
	if err != nil {
		return 0, nil
	}

	return *tokenEnable.Free, nil

}

func (e *Exchanges) GetPrices(cexName, token string) (float64, error) {
	cex, err := e.getCEX(cexName)
	if err != nil {
		return 0, nil
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

func (e *Exchanges) GetChains(cexName, token, withdrawChain string) (*models.ChainList, error) {
	cex, err := e.getCEX(cexName)
	if err != nil {
		return nil, err
	}

	return cex.GetChains(token, withdrawChain)
}

func (e *Exchanges) getCEX(cexName string) (ExchangeModule, error) {
	if cex, ok := e.CEXs[cexName]; ok {
		return cex, nil
	}
	return nil, fmt.Errorf("failed to get cex from adapters")
}
