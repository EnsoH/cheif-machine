package modules

func validateParam(param interface{}) bool {
	// Проверка, что это строка, и она не пуста. В дальнейшем переработать на универсальную функцию и возможно переместить в глобальные хелперы
	return param != ""
}
