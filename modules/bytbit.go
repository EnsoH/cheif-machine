package modules

import (
	"crypto/hmac"
	"crypto/sha256"
	"cw/logger"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type BybitModule struct {
	WithdrawEndpoint string
	API_key          string
	API_secret       string
}

func NewBybitModule(endpoint string, apiKey, apiSecret string) (*BybitModule, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, fmt.Errorf("Missing bybit api endpoint. Check config.")
	}

	return &BybitModule{
		WithdrawEndpoint: endpoint,
		API_key:          apiKey,
		API_secret:       apiSecret,
	}, nil
}

func (b *BybitModule) Action() error {
	exchange := ccxt.NewBybit(map[string]interface{}{
		"apiKey":          b.API_key,
		"secret":          b.API_secret,
		"enableRateLimit": true,
	})
	amount := 7
	network := "ARBI"
	addr := "0x99"
	// Подменяем HTTP-клиент в CCXT
	_, err := exchange.Withdraw(
		"USDT",
		7.0,
		"test_wallet",
		ccxt.WithWithdrawParams(map[string]interface{}{
			"forceChain": 1,
			"network":    "ARBI",
		}),
	)
	if err != nil {
		logger.GlobalLogger.Error(err)
		return err
	}

	logger.GlobalLogger.Infof("[%s] Withdraw %d successful. Chain %s", addr, amount, network)
	return nil
}

func (b *BybitModule) GetBalances() error {
	baseURL := "https://api.bybit.com"
	endpoint := "/v5/asset/transfer/query-account-coins-balance"
	// https: //api.bybit.com/v5/account/wallet-balance
	// Параметры запроса
	params := url.Values{}
	params.Add("accountType", "FUND")
	params.Add("coin", "USDT")

	// Создание метки времени
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Создание строки для подписи
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

func (b *BybitModule) GetPrices() error {
	baseURL := "https://api.bybit.com"
	endpoint := "/v5/market/tickers"

	params := url.Values{}
	params.Add("category", "spot") // "spot" - для спотового рынка, "linear" - для фьючерсов
	params.Add("symbol", "BTCUSDT")
	reqURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())

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
	// Создаем HTTP-запрос
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return err
	}

	// Отправка запроса
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return err
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return err
	}

	// Выводим ответ
	fmt.Println("Response:", string(body))
	return nil
}
