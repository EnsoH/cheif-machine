package modules

import "fmt"

func validateParam(param interface{}) bool {
	// Проверка, что это строка, и она не пуста. В дальнейшем переработать на универсальную функцию и возможно переместить в глобальные хелперы
	return param != ""
}

func createUrl(baseURL string, param ...string) string {
	url := fmt.Sprintf("%s?%s", baseURL, param[0])
	return url
}
