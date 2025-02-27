package modules

import (
	"crypto/hmac"
	"crypto/sha256"
	"cw/httpClient"
	"cw/logger"
	"cw/models"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
}

func NewBybitModule(balanceEndpoint, tickersEndpoint string, apiKey, apiSecret string, hc *httpClient.HttpClient) (*BybitModule, error) {
	if balanceEndpoint == "" || tickersEndpoint == "" {
		return nil, fmt.Errorf("Missing bybit api endpoint. Check config.")
	}

	return &BybitModule{
		BalanceEndpoint: balanceEndpoint,
		TickersEndpoint: tickersEndpoint,
		API_key:         apiKey,
		API_secret:      apiSecret,
		HttpClient:      hc,
	}, nil
}

func (b *BybitModule) Withdraw(token, address, network string, amount float64) error {
	exchange := ccxt.NewBybit(map[string]interface{}{
		"apiKey":          b.API_key,
		"secret":          b.API_secret,
		"enableRateLimit": true,
	})

	tx, err := exchange.Withdraw(
		token,
		amount,
		address,
		ccxt.WithWithdrawParams(map[string]interface{}{
			"forceChain": 1,
			"network":    network,
		}),
	)

	if err != nil {
		logger.GlobalLogger.Error(err)
		return err
	}

	logger.GlobalLogger.Infof("[%s] Withdraw %d successful. Chain %s. TxId: %v", address, amount, network, tx.TxId)
	return nil
}

func (b *BybitModule) GetBalances(token string) error {
	baseURL := "https://api.bybit.com"
	endpoint := "/v5/asset/transfer/query-account-coins-balance"

	// Параметры запроса
	params := url.Values{}
	params.Add("accountType", "FUND")
	params.Add("coin", token) // Замените на нужную валюту

	// Создание метки времени
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Формируем строку для подписи
	queryString := params.Encode()
	preSign := timestamp + b.API_key + "5000" + queryString
	// Создание HMAC-SHA256 подписи
	h := hmac.New(sha256.New, []byte(b.API_secret))
	h.Write([]byte(preSign))
	signature := hex.EncodeToString(h.Sum(nil))

	// Создание URL с параметрами
	reqURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, queryString)

	// Создание HTTP-запроса
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return err
	}

	// Добавление заголовков
	req.Header.Add("X-BAPI-SIGN", signature)
	req.Header.Add("X-BAPI-API-KEY", b.API_key)
	req.Header.Add("X-BAPI-TIMESTAMP", timestamp)
	req.Header.Add("X-BAPI-RECV-WINDOW", "5000")

	transport := &http.Transport{}
	proxy, err := url.Parse("")
	if err != nil {
		return err
	}
	transport.Proxy = http.ProxyURL(proxy)
	// Отправка запроса
	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return err
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return err
	}

	// Вывод ответа
	log.Printf("bal: %v", string(body))

	return nil
}

func (b *BybitModule) GetPrices(token string) error {
	url := fmt.Sprintf("%s?category=spot&symbol=%sUSDT", b.TickersEndpoint, token)

	var resp *models.BybitTickerResponse
	if err := b.HttpClient.SendJSONRequest(url, "GET", nil, &resp); err != nil {
		logger.GlobalLogger.Error(err)
		return err
	}
	logger.GlobalLogger.Infof("price: %v", resp.Result.List[0].LastPrice)
	return nil
}
