package httpClient

// func (h *HttpClient) setHeaders(req *http.Request) {
// 	params := url.Values{}
// 	params.Add("accountType", "FUND")
// 	params.Add("coin", token)

// 	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
// 	queryString := params.Encode()
// 	preSign := timestamp + b.API_key + "5000" + queryString
// 	// Создание HMAC-SHA256 подписи
// 	h = hmac.New(sha256.New, []byte(b.API_secret))
// 	h.Write([]byte(preSign))
// 	signature := hex.EncodeToString(h.Sum(nil))

// 	req.Header.Add("X-BAPI-SIGN", signature)
// 	req.Header.Add("X-BAPI-API-KEY", b.API_key)
// 	req.Header.Add("X-BAPI-TIMESTAMP", timestamp)
// 	req.Header.Add("X-BAPI-RECV-WINDOW", "5000")
// }

// func (h *HttpClient) getRandomUserAgent() string {
// r := rand.New(rand.NewSource(time.Now().UnixNano()))
// return globals.UserAgents[r.Intn(len(globals.UserAgents))]
// }

// func (h *HttpClient) getSecChUa(userAgent string) (string, string) {
// 	if strings.Contains(userAgent, "Macintosh") {
// 		return globals.SecChUa["Macintosh"], globals.Platforms["Macintosh"]
// 	} else if strings.Contains(userAgent, "Windows") {
// 		return globals.SecChUa["Windows"], globals.Platforms["Windows"]
// 	} else if strings.Contains(userAgent, "Linux") {
// 		return globals.SecChUa["Linux"], globals.Platforms["Linux"]
// 	}
// 	return globals.SecChUa["Unknown"], `"Unknown"`
// }
