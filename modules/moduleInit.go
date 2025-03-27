package modules

import (
	"cw/account"
	"cw/config"
	"cw/httpClient"
	"cw/models"
	"fmt"
	"math/big"
	"strings"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

// Interface for Exchange Adapters.
type ExchangeModule interface {
	Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error)
	GetBalance(symbol string) (ccxt.Balances, error)
	GetPrices(token string) (float64, error)
	GetChains(token, withdrawChain string) (*models.ChainList, error)
}

// Interface for Bridge Adapters.
type BridgerModule interface {
	Bridge(fromChain, destChain, fromToken, toToken string, amount *big.Int, acc *account.Account) error
}

// Interface for wallet generatos
type WallGenModule interface {
	GenerateWallet() (privateKey string, address string, mnemonic string, err error)
}

type WarmerModule interface {
	Approve() error
	SlefTransfer() error
	Wrap() error
}

type Modules struct {
	Exchange         *Exchanges
	Bridges          *Bridger
	WalletGenerators *WalletGenerator
	Collector        *Collectors
}

// Default init modules.
// TODO: create 'type' for modules. Delete hard code 'module'
func ModulesInit(moduleChoice string, names ...string) (*Modules, error) {
	mod := Modules{}

	switch moduleChoice {
	case "CexWithdrawer":
		exchange, err := ExchangeInit(names...)
		if err != nil {
			return nil, err
		}
		if exchange == nil {
			return nil, fmt.Errorf("не было инициализировано ни одной биржи, проверьте биржевые конфиги на наличие всех параметров модулей: [%s]", strings.Join(names, ", "))
		}

		mod.Exchange = exchange
	case "Bridger":
		bridger, err := BridgeInit(names...)
		if err != nil {
			return nil, err
		}

		if bridger == nil {
			return nil, fmt.Errorf("не было инициализировано ни одного моста, проверьте биржевые конфиги на наличие всех параметров модулей: [%s]", strings.Join(names, ", "))
		}
		mod.Bridges = bridger
	case "Cex_Bridger":
		exchange, err := ExchangeInit(names...)
		if err != nil {
			return nil, err
		}
		if exchange == nil {
			return nil, fmt.Errorf("не было инициализировано ни одной биржи, проверьте биржевые конфиги на наличие всех параметров модулей: [%s]", strings.Join(names, ", "))
		}

		bridger, err := BridgeInit(names...)
		if err != nil {
			return nil, err
		}
		if bridger == nil {
			return nil, fmt.Errorf("не было инициализировано ни одного моста, проверьте биржевые конфиги на наличие всех параметров модулей: [%s]", strings.Join(names, ", "))
		}

		mod.Exchange = exchange
		mod.Bridges = bridger
	case "WalletGenerator":
		generator, err := NewWalletGen(config.UserCfg.WalletGeneratorCfg.WalletType)
		if err != nil {
			return nil, err
		}
		mod.WalletGenerators = generator
	case "Сollector":
		collector := NewCollector()
		mod.Collector = collector
	default:
		return nil, fmt.Errorf("неизвестный модуль: %s", moduleChoice)
	}

	return &mod, nil
}

func ExchangeInit(cexNames ...string) (*Exchanges, error) {
	hc, err := httpClient.NewHttpClient(
		httpClient.WithHttp2(),
		httpClient.WithProxy(config.Cfg.IpAddresses[0]),
	)
	if err != nil {
		return nil, err
	}

	return NewExchangeModule(hc, cexNames...)
}

// ДОБАВИТЬ ПРОВЕРКУ НА ВАЛИДНОСТЬ ПРОКСИ!
func BridgeInit(bridgeNames ...string) (*Bridger, error) {
	hc, err := httpClient.NewHttpClient(
		httpClient.WithHttp2(),
		httpClient.WithProxy(config.Cfg.IpAddresses[0]),
	)
	if err != nil {
		return nil, err
	}

	return NewBridgeModule(hc, bridgeNames...)
}
