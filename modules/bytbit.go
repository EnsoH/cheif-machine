package modules

import (
	"crypto/hmac"
	"crypto/sha256"
	"cw/config"
	"cw/httpClient"
	"cw/logger"
	"cw/models"
	"cw/utils"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type BybitModule struct {
	BalanceEndpoint string
	TickersEndpoint string
	API_key         string
	API_secret      string
	HttpClient      *httpClient.HttpClient
	BybitEx         *ccxt.Bybit
}

func NewBybitModule(balanceEndpoint, tickersEndpoint, apiKey, apiSecret, proxy string, hc *httpClient.HttpClient) (*BybitModule, error) {
	if balanceEndpoint == "" || tickersEndpoint == "" {
		return nil, fmt.Errorf("missing bybit api endpoint, check config")
	}

	bybit := ccxt.NewBybit(map[string]interface{}{
		"apiKey":          apiKey,
		"secret":          apiSecret,
		"enableRateLimit": true,
		"proxy":           config.Cfg.IpAddresses[0],
	})
	bybit.HttpProxy = proxy

	return &BybitModule{
		BalanceEndpoint: balanceEndpoint,
		TickersEndpoint: tickersEndpoint,
		API_key:         apiKey,
		API_secret:      apiSecret,
		HttpClient:      hc,
		BybitEx:         &bybit,
	}, nil
}

func (b *BybitModule) Withdraw(token, address, network string, amount float64) error {
	str := strconv.FormatFloat(amount, 'f', 6, 64)
	result, _ := strconv.ParseFloat(str, 64)

	tx, err := b.BybitEx.Withdraw(
		token,
		result,
		address,
		ccxt.WithdrawOptions(ccxt.WithWithdrawParams(map[string]interface{}{
			"forceChain": 1,
			"network":    network,
		})),
	)

	if err != nil {
		log.Println("Ошибка при выводе средств:", err)
		return err
	}

	logger.GlobalLogger.Infof("[%s] Withdraw %s успешен. Chain %s. Amount %f TxId: %v", address, token, network, amount, tx.TxId)
	return nil
}

func (b *BybitModule) GetBalances(token string) (float64, error) {
	reqURL := fmt.Sprintf(b.BalanceEndpoint, token)
	headers := createHeaders(token, b.API_key, b.API_secret)

	var resp *models.BybitBalanceResponse
	if err := b.HttpClient.SendJSONRequest(reqURL, "GET", nil, &resp, headers); err != nil {
		return 0.0, err
	}

	return utils.СonvertStringToFloat(resp.Result.Balance[0].WalletBalance)
}

func (b *BybitModule) GetPrices(token string) (float64, error) {
	if token == "USDT" || token == "USDC" {
		return 1.0, nil
	}

	var resp *models.BybitTickerResponse
	if err := b.HttpClient.SendJSONRequest(fmt.Sprintf("%s?category=spot&symbol=%sUSDT", b.TickersEndpoint, token), "GET", nil, &resp, map[string]string{}); err != nil {
		logger.GlobalLogger.Error(err)
		return 0.0, err
	}

	return utils.СonvertStringToFloat(resp.Result.List[0].LastPrice)
}

func createHeaders(token, api_key, api_secret string) map[string]string {
	params := url.Values{}
	params.Add("accountType", "FUND")
	params.Add("coin", token)
	queryString := params.Encode()

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	preSign := timestamp + api_key + "5000" + queryString
	h := hmac.New(sha256.New, []byte(api_secret))
	h.Write([]byte(preSign))
	signature := hex.EncodeToString(h.Sum(nil))

	headers := map[string]string{
		"X-BAPI-SIGN":        signature,
		"X-BAPI-API-KEY":     api_key,
		"X-BAPI-TIMESTAMP":   timestamp,
		"X-BAPI-RECV-WINDOW": "5000",
	}

	return headers
}
