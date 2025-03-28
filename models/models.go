package models

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ##############################
// ######### APP CONFIG #########
// ##############################
type AppConfig struct {
	Threads            int               `json:"threads"`
	IpAddresses        []string          `json:"ip_addresses"`
	CEXConfigs         CexConfig         `json:"cex"`
	Rpc                map[string]string `json:"rpc"`
	Endpoints          map[string]string `json:"endpoints"`
	AttentionGwei      string            `json:"attention_gwei"`
	AttentionTimeCycle int               `json:"attention_time_cycle"`
	AttentionMaxTime   int               `json:"max_attention_time"`
	DepositWaitingTime int               `json:"deposit_waiting_time"`
	SleepAfterWithdraw int               `json:"sleep_after_withdraw"`
}

type CexConfig struct {
	BybitCfg   Bybit   `json:"bybit"`
	BinanceCfg Binance `json:"binance"`
	MexcCfg    Mexc    `json:"mexc"`
	KucoinCfg  Kucoin  `json:"kucoin"`
	OkxCfg     Okx     `json:"okx"`
}

// ##############################
// ######### USER CONFIG ########
// ##############################
type UserConfig struct {
	WithdrawConfig     WithdrawConfig     `json:"cex_withdraw"`
	BridgeConfig       BridgeConfig       `json:"bridge"`
	WalletGeneratorCfg WalletGeneratorCfg `json:"wallet_generator"`
	CollectorConfig    CollectorConfig    `json:"collector_config"`
}

// ##############################
// ######### CEX CONFIG ########
// ##############################
type Bybit struct {
	API_key    string `json:"api_key"`
	API_secret string `json:"secret_key"`
}

type Binance struct {
	API_key    string `json:"api_key"`
	API_secret string `json:"secret_key"`
}

type Mexc struct {
	API_key    string `json:"api_key"`
	API_secret string `json:"secret_key"`
}

type Kucoin struct {
	API_key    string `json:"api_key"`
	API_secret string `json:"secret_key"`
	Password   string `json:"password"`
}

type Okx struct {
	API_key    string `json:"api_key"`
	API_secret string `json:"secret_key"`
	Password   string `json:"password"`
}

// #########################################
// ############ WITHDRAW OPTIONS ###########
// #########################################
type WithdrawConfig struct {
	CEX                 string    `json:"cex"`
	Chain               []string  `json:"chain"`
	Currency            []string  `json:"currency"`
	AmountRange         []float64 `json:"amount_range"`
	TimeRange           []float64 `json:"time_range"`
	DestinationChain    string    `json:"destination_chain"`
	DestinationCurrency string    `json:"destination_currency"`
}

type WithdrawAction struct {
	CEX                 string
	Address             string
	Chain               string
	Amount              float64
	Currency            string
	TimeRange           float64
	DestinationChain    string
	DestinationCurrency string
}

type ChainList struct {
	Chain       string
	Network     string
	WithdrawFee float64
}

type TokenInfo struct {
	AvailableCapacity float64  // Доступная мощность в USD
	Price             float64  // Цена токена
	Chains            []string // Список доступных сетей
	Balance           float64  // Баланс токена
}

// #########################################
// ############ BRIDGE OPTIONS ###########
// #########################################
type BridgeConfig struct {
	Bridge       string    `json:"bridge"`
	FromChain    string    `json:"from_chain"`
	ToChain      string    `json:"to_chain"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	AmountRange  []float64 `json:"amount_range"`
	TimeRange    []float64 `json:"time_range"`
}

type BridgeAction struct {
	Address      common.Address
	Bridge       string
	FromChain    string
	ToChain      string
	FromCurrency string
	ToCurrency   string
	Amount       *big.Int
	Time         float64
}

// #########################################
// ######## WALLET GENERATOR CONFIF ########
// #########################################

type WalletGeneratorCfg struct {
	WalletType  string   `json:"wallet_type"`
	WalletCount uint32   `json:"wallet_count"`
	CsvHeaders  []string `json:"csv_headers"`
}

type CollectorConfig struct {
	DestinationAddress   string            `json:"destination_address"`
	DestinationAddresses map[string]string `json:"destination_addresses"`
	Chains               []string          `json:"chains"`
}

// #################################
// ####### RESPONSE SHEMES #########
// #################################

type BybitCurrencyList struct {
	Active  bool    `json:"active"`
	Code    string  `json:"code"`
	Deposit bool    `json:"deposit"`
	Fee     float32 `json:"fee"`
	Id      string  `json:"id"`
	Info    struct {
		Coin         string          `json:"coin"`
		Name         string          `json:"name"`
		RemainAmount json.RawMessage `json:"remainAmount"` // Может быть строкой или числом
		Chains       []struct {
			Chain           string          `json:"chain"`
			ChainType       string          `json:"chainType"`
			ContractAddress string          `json:"contractAddress"`
			Confirmation    json.RawMessage `json:"confirmation"` // Может быть строкой или числом
			DepositMin      json.RawMessage `json:"depositMin"`
			WithdrawFee     json.RawMessage `json:"withdrawFee"`
			WithdrawMin     json.RawMessage `json:"withdrawMin"`
		} `json:"chains"`
	} `json:"info"`
}

type BinanceCurrencyList struct {
	Active  bool               `json:"active"`
	Code    string             `json:"code"`
	Deposit bool               `json:"deposit"`
	Fee     float64            `json:"fee"`
	Fees    map[string]float64 `json:"fees"`
	Id      string             `json:"id"`
	Info    struct {
		Coin             string          `json:"coin"`
		DepositAllEnable bool            `json:"depositAllEnable"`
		Free             json.RawMessage `json:"free"`
		Freeze           json.RawMessage `json:"freeze"`
		Ipoable          json.RawMessage `json:"ipoable"`
		Ipoing           json.RawMessage `json:"ipoing"`
		IsLegalMoney     bool            `json:"isLegalMoney"`
		Locked           json.RawMessage `json:"locked"`
		Name             string          `json:"name"`
		NetworkList      []struct {
			AddressRegex            string          `json:"addressRegex"`
			Busy                    bool            `json:"busy"`
			Coin                    string          `json:"coin"`
			ContractAddress         string          `json:"contractAddress,omitempty"`
			ContractAddressUrl      string          `json:"contractAddressUrl"`
			DepositDesc             string          `json:"depositDesc"`
			DepositDust             json.RawMessage `json:"depositDust"`
			DepositEnable           bool            `json:"depositEnable"`
			EstimatedArrivalTime    json.RawMessage `json:"estimatedArrivalTime"`
			IsDefault               bool            `json:"isDefault"`
			MemoRegex               string          `json:"memoRegex"`
			MinConfirm              json.RawMessage `json:"minConfirm"`
			Name                    string          `json:"name"`
			Network                 string          `json:"network"`
			ResetAddressStatus      bool            `json:"resetAddressStatus"`
			SameAddress             bool            `json:"sameAddress"`
			SpecialTips             string          `json:"specialTips"`
			SpecialWithdrawTips     string          `json:"specialWithdrawTips"`
			UnLockConfirm           json.RawMessage `json:"unLockConfirm"`
			WithdrawDesc            string          `json:"withdrawDesc"`
			WithdrawFee             json.RawMessage `json:"withdrawFee"`
			WithdrawIntegerMultiple json.RawMessage `json:"withdrawIntegerMultiple"`
			WithdrawInternalMin     json.RawMessage `json:"withdrawInternalMin"`
			WithdrawMax             json.RawMessage `json:"withdrawMax"`
			WithdrawMin             json.RawMessage `json:"withdrawMin"`
		} `json:"networkList"`
	} `json:"info"`
	// Если структура ограничений известна, можно заменить json.RawMessage на конкретную структуру
	Limits   json.RawMessage `json:"limits"`
	Margin   interface{}     `json:"margin"`
	Name     string          `json:"name"`
	Networks map[string]struct {
		Active  bool            `json:"active"`
		Deposit bool            `json:"deposit"`
		Fee     json.RawMessage `json:"fee"`
		Id      string          `json:"id"`
		Info    struct {
			AddressRegex            string          `json:"addressRegex"`
			Busy                    bool            `json:"busy"`
			Coin                    string          `json:"coin"`
			ContractAddress         string          `json:"contractAddress,omitempty"`
			ContractAddressUrl      string          `json:"contractAddressUrl"`
			DepositDesc             string          `json:"depositDesc"`
			DepositDust             json.RawMessage `json:"depositDust"`
			DepositEnable           bool            `json:"depositEnable"`
			EstimatedArrivalTime    json.RawMessage `json:"estimatedArrivalTime"`
			IsDefault               bool            `json:"isDefault"`
			MemoRegex               string          `json:"memoRegex"`
			MinConfirm              json.RawMessage `json:"minConfirm"`
			Name                    string          `json:"name"`
			Network                 string          `json:"network"`
			ResetAddressStatus      bool            `json:"resetAddressStatus"`
			SameAddress             bool            `json:"sameAddress"`
			SpecialTips             string          `json:"specialTips"`
			SpecialWithdrawTips     string          `json:"specialWithdrawTips"`
			UnLockConfirm           json.RawMessage `json:"unLockConfirm"`
			WithdrawDesc            string          `json:"withdrawDesc"`
			WithdrawFee             json.RawMessage `json:"withdrawFee"`
			WithdrawIntegerMultiple json.RawMessage `json:"withdrawIntegerMultiple"`
			WithdrawInternalMin     json.RawMessage `json:"withdrawInternalMin"`
			WithdrawMax             json.RawMessage `json:"withdrawMax"`
			WithdrawMin             json.RawMessage `json:"withdrawMin"`
		} `json:"info"`
		Limits struct {
			Deposit  map[string]json.RawMessage `json:"deposit"`
			Withdraw map[string]json.RawMessage `json:"withdraw"`
		} `json:"limits"`
		Network   string          `json:"network"`
		Precision json.RawMessage `json:"precision"`
		Withdraw  bool            `json:"withdraw"`
	} `json:"networks"`
	Precision json.RawMessage `json:"precision"`
	Type      string          `json:"type"`
	Withdraw  bool            `json:"withdraw"`
}

type MexcCurrencyList struct {
	Active  bool    `json:"active"`
	Code    string  `json:"code"`
	Deposit bool    `json:"deposit"`
	Fee     float64 `json:"fee"`
	Id      string  `json:"id"`
	Info    struct {
		Coin        string `json:"coin"`
		Name        string `json:"name"`
		NetworkList []struct {
			Coin                    string          `json:"coin"`
			Contract                string          `json:"contract"`
			DepositDesc             *string         `json:"depositDesc"` // может быть null
			DepositEnable           bool            `json:"depositEnable"`
			DepositTips             *string         `json:"depositTips"` // может быть null или пустой строкой
			MinConfirm              int             `json:"minConfirm"`
			Name                    string          `json:"name"`
			NetWork                 string          `json:"netWork"`
			Network                 string          `json:"network"`
			SameAddress             bool            `json:"sameAddress"`
			WithdrawEnable          bool            `json:"withdrawEnable"`
			WithdrawFee             json.RawMessage `json:"withdrawFee"`
			WithdrawIntegerMultiple json.RawMessage `json:"withdrawIntegerMultiple"`
			WithdrawMax             json.RawMessage `json:"withdrawMax"`
			WithdrawMin             json.RawMessage `json:"withdrawMin"`
			WithdrawTips            *string         `json:"withdrawTips"`
		} `json:"networkList"`
	} `json:"info"`
	Limits struct {
		Amount struct {
			Max *float64 `json:"max"` // может быть null
			Min *float64 `json:"min"` // может быть null
		} `json:"amount"`
		Withdraw struct {
			Max json.RawMessage `json:"max"`
			Min json.RawMessage `json:"min"`
		} `json:"withdraw"`
	} `json:"limits"`
	Name     string `json:"name"`
	Networks map[string]struct {
		Active  bool    `json:"active"`
		Deposit bool    `json:"deposit"`
		Fee     float64 `json:"fee"`
		Id      string  `json:"id"`
		Info    struct {
			Coin                    string          `json:"coin"`
			Contract                string          `json:"contract"`
			DepositDesc             *string         `json:"depositDesc"`
			DepositEnable           bool            `json:"depositEnable"`
			DepositTips             *string         `json:"depositTips"`
			MinConfirm              int             `json:"minConfirm"`
			Name                    string          `json:"name"`
			NetWork                 string          `json:"netWork"`
			Network                 string          `json:"network"`
			SameAddress             bool            `json:"sameAddress"`
			WithdrawEnable          bool            `json:"withdrawEnable"`
			WithdrawFee             json.RawMessage `json:"withdrawFee"`
			WithdrawIntegerMultiple json.RawMessage `json:"withdrawIntegerMultiple"`
			WithdrawMax             json.RawMessage `json:"withdrawMax"`
			WithdrawMin             json.RawMessage `json:"withdrawMin"`
			WithdrawTips            *string         `json:"withdrawTips"`
		} `json:"info"`
		Limits struct {
			Withdraw struct {
				Max json.RawMessage `json:"max"`
				Min json.RawMessage `json:"min"`
			} `json:"withdraw"`
		} `json:"limits"`
		Network   string          `json:"network"`
		Precision json.RawMessage `json:"precision"`
		Withdraw  bool            `json:"withdraw"`
	} `json:"networks"`
	Precision json.RawMessage `json:"precision"`
	Withdraw  bool            `json:"withdraw"`
}

type KucoinCurrencyList struct {
	Active  bool                       `json:"active"`
	Code    string                     `json:"code"`
	Deposit bool                       `json:"deposit"`
	Fee     float64                    `json:"fee"`
	Fees    map[string]json.RawMessage `json:"fees"`
	Id      string                     `json:"id"`
	Info    struct {
		Chains []struct {
			ChainId           string          `json:"chainId"`
			ChainName         string          `json:"chainName"`
			Confirms          int             `json:"confirms"`
			ContractAddress   json.RawMessage `json:"contractAddress"`
			DepositMinSize    json.RawMessage `json:"depositMinSize"`
			IsDepositEnabled  bool            `json:"isDepositEnabled"`
			IsWithdrawEnabled bool            `json:"isWithdrawEnabled"`
			MaxDeposit        json.RawMessage `json:"maxDeposit"`
			MaxWithdraw       json.RawMessage `json:"maxWithdraw"`
			NeedTag           bool            `json:"needTag"`
			PreConfirms       int             `json:"preConfirms"`
			WithdrawFeeRate   json.RawMessage `json:"withdrawFeeRate"`
			WithdrawPrecision int             `json:"withdrawPrecision"`
			WithdrawalMinFee  json.RawMessage `json:"withdrawalMinFee"`
			WithdrawalMinSize json.RawMessage `json:"withdrawalMinSize"`
		} `json:"chains"`
		Confirms        json.RawMessage `json:"confirms"`
		ContractAddress json.RawMessage `json:"contractAddress"`
		Currency        string          `json:"currency"`
		FullName        string          `json:"fullName"`
		IsDebitEnabled  bool            `json:"isDebitEnabled"`
		IsMarginEnabled bool            `json:"isMarginEnabled"`
		Name            string          `json:"name"`
		Precision       json.RawMessage `json:"precision"`
	} `json:"info"`
	Limits struct {
		Deposit struct {
			Max json.RawMessage `json:"max"`
			Min float64         `json:"min"`
		} `json:"deposit"`
		Withdraw struct {
			Max json.RawMessage `json:"max"`
			Min float64         `json:"min"`
		} `json:"withdraw"`
	} `json:"limits"`
	Name     string `json:"name"`
	Networks map[string]struct {
		Active  bool    `json:"active"`
		Code    string  `json:"code"`
		Deposit bool    `json:"deposit"`
		Fee     float64 `json:"fee"`
		Id      string  `json:"id"`
		Info    struct {
			ChainId           string          `json:"chainId"`
			ChainName         string          `json:"chainName"`
			Confirms          int             `json:"confirms"`
			ContractAddress   json.RawMessage `json:"contractAddress"`
			DepositMinSize    json.RawMessage `json:"depositMinSize"`
			IsDepositEnabled  bool            `json:"isDepositEnabled"`
			IsWithdrawEnabled bool            `json:"isWithdrawEnabled"`
			MaxDeposit        json.RawMessage `json:"maxDeposit"`
			MaxWithdraw       json.RawMessage `json:"maxWithdraw"`
			NeedTag           bool            `json:"needTag"`
			PreConfirms       int             `json:"preConfirms"`
			WithdrawFeeRate   json.RawMessage `json:"withdrawFeeRate"`
			WithdrawPrecision int             `json:"withdrawPrecision"`
			WithdrawalMinFee  json.RawMessage `json:"withdrawalMinFee"`
			WithdrawalMinSize json.RawMessage `json:"withdrawalMinSize"`
		} `json:"info"`
		Limits struct {
			Deposit struct {
				Max json.RawMessage `json:"max"`
				Min json.RawMessage `json:"min"`
			} `json:"deposit"`
			Withdraw struct {
				Max json.RawMessage `json:"max"`
				Min json.RawMessage `json:"min"`
			} `json:"withdraw"`
		} `json:"limits"`
		Name      string          `json:"name"`
		Precision json.RawMessage `json:"precision"`
		Withdraw  bool            `json:"withdraw"`
	} `json:"networks"`
	NumericId json.RawMessage `json:"numericId"`
	Precision json.RawMessage `json:"precision"`
	Type      string          `json:"type"`
	Withdraw  bool            `json:"withdraw"`
}

type OkxCurrenceList struct {
	Active  bool     `json:"active"`
	Code    string   `json:"code"`
	Deposit bool     `json:"deposit"`
	Fee     *float64 `json:"fee"` // может быть nil
	Id      string   `json:"id"`
	Info    []struct {
		BurningFeeRate       json.RawMessage `json:"burningFeeRate"` // может быть nil
		CanDep               bool            `json:"canDep"`
		CanInternal          bool            `json:"canInternal"`
		CanWd                bool            `json:"canWd"`
		Ccy                  string          `json:"ccy"`
		Chain                string          `json:"chain"`
		CtAddr               string          `json:"ctAddr"`
		DepEstOpenTime       json.RawMessage `json:"depEstOpenTime"`      // метка времени или nil
		DepQuotaFixed        json.RawMessage `json:"depQuotaFixed"`       // может быть nil
		DepQuoteDailyLayer2  json.RawMessage `json:"depQuoteDailyLayer2"` // может быть nil
		Fee                  json.RawMessage `json:"fee"`
		LogoLink             string          `json:"logoLink"`
		MainNet              bool            `json:"mainNet"`
		MaxFee               json.RawMessage `json:"maxFee"`
		MaxFeeForCtAddr      json.RawMessage `json:"maxFeeForCtAddr"` // может быть nil
		MaxWd                json.RawMessage `json:"maxWd"`
		MinDep               json.RawMessage `json:"minDep"`
		MinDepArrivalConfirm json.RawMessage `json:"minDepArrivalConfirm"`
		MinFee               json.RawMessage `json:"minFee"`
		MinFeeForCtAddr      json.RawMessage `json:"minFeeForCtAddr"` // может быть nil
		MinInternal          json.RawMessage `json:"minInternal"`
		MinWd                interface{}     `json:"minWd"`              // может быть числом или строкой
		MinWdUnlockConfirm   interface{}     `json:"minWdUnlockConfirm"` // может быть числом или строкой
		Name                 string          `json:"name"`
		NeedTag              bool            `json:"needTag"`
		UsedDepQuotaFixed    interface{}     `json:"usedDepQuotaFixed"` // может быть числом или строкой
		UsedWdQuota          interface{}     `json:"usedWdQuota"`       // может быть числом или строкой
		WdEstOpenTime        json.RawMessage `json:"wdEstOpenTime"`     // может быть nil
		WdQuota              json.RawMessage `json:"wdQuota"`
		WdTickSz             json.RawMessage `json:"wdTickSz"`
	} `json:"info"`
	Limits struct {
		Amount struct {
			Max interface{} `json:"max"`
			Min interface{} `json:"min"`
		} `json:"amount"`
	} `json:"limits"`
	Name     string `json:"name"`
	Networks map[string]struct {
		Active  bool    `json:"active"`
		Deposit bool    `json:"deposit"`
		Fee     float64 `json:"fee"`
		Id      string  `json:"id"`
		Info    struct {
			BurningFeeRate       json.RawMessage `json:"burningFeeRate"` // может быть nil
			CanDep               bool            `json:"canDep"`
			CanInternal          bool            `json:"canInternal"`
			CanWd                bool            `json:"canWd"`
			Ccy                  string          `json:"ccy"`
			Chain                string          `json:"chain"`
			CtAddr               string          `json:"ctAddr"`
			DepEstOpenTime       json.RawMessage `json:"depEstOpenTime"`      // метка времени или nil
			DepQuotaFixed        json.RawMessage `json:"depQuotaFixed"`       // может быть nil
			DepQuoteDailyLayer2  json.RawMessage `json:"depQuoteDailyLayer2"` // может быть nil
			Fee                  json.RawMessage `json:"fee"`
			LogoLink             string          `json:"logoLink"`
			MainNet              bool            `json:"mainNet"`
			MaxFee               json.RawMessage `json:"maxFee"`
			MaxFeeForCtAddr      json.RawMessage `json:"maxFeeForCtAddr"` // может быть nil
			MaxWd                json.RawMessage `json:"maxWd"`
			MinDep               json.RawMessage `json:"minDep"`
			MinDepArrivalConfirm json.RawMessage `json:"minDepArrivalConfirm"`
			MinFee               json.RawMessage `json:"minFee"`
			MinFeeForCtAddr      json.RawMessage `json:"minFeeForCtAddr"` // может быть nil
			MinInternal          json.RawMessage `json:"minInternal"`
			MinWd                interface{}     `json:"minWd"`              // может быть числом или строкой
			MinWdUnlockConfirm   interface{}     `json:"minWdUnlockConfirm"` // может быть числом или строкой
			Name                 string          `json:"name"`
			NeedTag              bool            `json:"needTag"`
			UsedDepQuotaFixed    interface{}     `json:"usedDepQuotaFixed"` // может быть числом или строкой
			UsedWdQuota          interface{}     `json:"usedWdQuota"`       // может быть числом или строкой
			WdEstOpenTime        json.RawMessage `json:"wdEstOpenTime"`     // может быть nil
			WdQuota              json.RawMessage `json:"wdQuota"`
			WdTickSz             json.RawMessage `json:"wdTickSz"`
		} `json:"info"`
		Limits struct {
			Withdraw struct {
				Max float64 `json:"max"`
				Min float64 `json:"min"`
			} `json:"withdraw"`
		} `json:"limits"`
		Network   string  `json:"network"`
		Precision float64 `json:"precision"`
		Withdraw  bool    `json:"withdraw"`
	} `json:"networks"`
	Precision float64 `json:"precision"`
	Withdraw  bool    `json:"withdraw"`
}

// #################################
// ######## Relay responses ########
// #################################

type RelayQuoteModel struct {
	TokenIn            string
	TokenOut           string
	OriginChainId      int64
	DestinationChainId int64
}

type RelayRequest struct {
	User                 string `json:"user"`
	OriginChainId        int64  `json:"originChainId"`
	DestinationChainId   int64  `json:"destinationChainId"`
	OriginCurrency       string `json:"originCurrency"`
	DestinationCurrency  string `json:"destinationCurrency"`
	Recipient            string `json:"recipient"`
	TradeType            string `json:"tradeType"`
	Amount               string `json:"amount"`
	Referrer             string `json:"referrer"`
	UseExternalLiquidity bool   `json:"useExternalLiquidity"`
	UseDepositAddress    bool   `json:"useDepositAddress"`
}

type RelayResponse struct {
	Steps []struct {
		Items []struct {
			Data struct {
				From    string `json:"from"`
				To      string `json:"to"`
				Data    string `json:"data"`
				Value   string `json:"value"`
				ChainID int    `json:"chainId"`
			} `json:"data"`
		} `json:"items"`
	} `json:"steps"`
	Details struct {
		Impact struct {
			Usd     string `json:"usd"`
			Percent string `json:"percent"`
		} `json:"totalImpact"`
		CurrencyIn struct {
			Currency struct {
				Address string `json:"address"`
			} `json:"currency"`
		} `json:"currencyIn"`
		CurrencyOut struct {
			Currency struct {
				Address string `json:"address"`
			} `json:"currency"`
		} `json:"currencyOut"`
	} `json:"details"`
	Message string `json:"message"`
	ErrCode string `json:"errorCode"`
}

type RelayTransactionData struct {
	TokenIn common.Address
	To      common.Address
	Value   *big.Int
	Data    []byte
}

// #################################
// ######## Cryptorank responses ########
// #################################
type CryptoRankPrice struct {
	Data []struct {
		Last float64 `json:"last"`
	} `json:"data"`
}
