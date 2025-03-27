package process

import (
	"cw/account"
	"cw/config"
	"cw/ethClient"
	"cw/globals"
	"cw/httpClient"
	"cw/logger"
	"cw/models"
	"cw/modules"
	"cw/utils"
	"fmt"
	"math"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
)

// ############## Bridge Helpers #################
func buildAccountMap(accounts []*account.Account) map[common.Address]*account.Account {
	accountMap := make(map[common.Address]*account.Account, len(accounts))
	for _, acc := range accounts {
		accountMap[acc.Address] = acc
	}
	return accountMap
}

func getTokenInfo(chain, currency string) (common.Address, *ethClient.EthClient, uint8, float64, error) {
	ethClientInstance := ethClient.GlobalETHClient[chain]

	tokenContract, err := getTokenContract(chain, currency)
	if err != nil {
		return common.Address{}, nil, 0, 0, err
	}

	decimals, err := ethClientInstance.GetDecimals(tokenContract)
	if err != nil {
		return common.Address{}, nil, 0, 0, err
	}

	tokenPrice, err := getTokenPrice(currency)
	if err != nil {
		return common.Address{}, nil, 0, 0, err
	}

	return tokenContract, ethClientInstance, decimals, tokenPrice, nil
}

func getAmount(chain, currency string, amountUsd float64) (*big.Int, error) {
	_, _, decimals, tokenPrice, err := getTokenInfo(chain, currency)
	if err != nil {
		return nil, err
	}

	amount := (amountUsd * 0.98) / tokenPrice
	weiAmount, err := utils.ConvertToWei(amount, int(decimals))
	if err != nil {
		return nil, err
	}

	return weiAmount, nil
}

func getTokenContract(chain, currency string) (common.Address, error) {
	chainId, err := ethClient.GlobalETHClient[chain].GetChainID()
	if err != nil {
		return common.Address{}, err
	}

	tokenContract, ok := globals.TokenContracts[chainId][currency]
	if !ok {
		return common.Address{}, fmt.Errorf("нет контракта токена")
	}

	return common.HexToAddress(tokenContract), nil
}

// ##################### CRYPTORANK HELPERS #########################
// getTokenPrice получает цену токена по его символу.
func getTokenPrice(tokenSymbol string) (float64, error) {
	client, err := httpClient.NewHttpClient()
	if err != nil {
		return 0, fmt.Errorf("не удалось создать http клиент: %w", err)
	}

	tokenName, ok := globals.TokenSymbolToName[tokenSymbol]
	if !ok {
		return 0, fmt.Errorf("токен %s не найден в конфигурации", tokenSymbol)
	}

	var priceData models.CryptoRankPrice
	url := fmt.Sprintf(config.Cfg.Endpoints["cryptorank"], tokenName)
	if err := client.SendJSONRequest(url, "GET", nil, &priceData, nil); err != nil {
		return 0, err
	}
	if len(priceData.Data) == 0 {
		return 0, fmt.Errorf("не получены данные цены для токена %s", tokenSymbol)
	}
	return priceData.Data[0].Last, nil
}

// #################### CEX WITHDRAW HELPERS ################
func getRandomChain(chains []string) string {
	return chains[rand.Intn(len(chains))]
}

func calculateAmount(token string, amount float64, exchange modules.Exchanges, cex string) (float64, error) {
	tickerPrice, err := exchange.GetPrices(cex, token)
	if err != nil {
		return 0.0, err
	}

	return (amount / tickerPrice), nil
}

// #################### GENERAL HELPERS ######################
func getRandomAmount(amountArr []float64) float64 {
	switch len(amountArr) {
	case 0:
		return 0
	case 1:
		return amountArr[0]
	default:
		min, max := amountArr[0], amountArr[1]
		if min > max {
			min, max = max, min
		}

		if min == max {
			return min
		}

		randoValue := min + rand.Float64()*(max-min)
		return math.Round(randoValue*100) / 100
	}
}

func handleWithdrawError(err error) error {
	errorContext, critical := utils.IsCriticalError(err)
	if critical {
		return err
	}
	logger.GlobalLogger.Warn("ошибка: %s, дальнейшее выполнение под вопросом", errorContext)
	return nil
}

// #################### LOGGER HELPERS #######################
func loggingBridgeAction(actions []*models.BridgeAction) {
	logger.GlobalLogger.Infof("################ BRIDGE INFO ################")
	logger.GlobalLogger.Infof("###############################################")
	for _, act := range actions {
		logger.GlobalLogger.Infof("[%s] %s | %s | %s | %s | %v", act.Address, act.FromChain, act.ToChain, act.FromCurrency, act.ToCurrency, act.Amount)
	}
}

func loggingActions(actions []models.WithdrawAction) {
	logger.GlobalLogger.Infof("################ WITHDRAW INFO ################")
	logger.GlobalLogger.Infof("###############################################")
	for _, act := range actions {
		logger.GlobalLogger.Infof("[%s] | %s | %s | %.8f|", act.Address, act.Chain, act.Currency, act.Amount)
	}
}
