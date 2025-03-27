package utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

func СonvertStringToFloat(line string) (float64, error) {
	return strconv.ParseFloat(line, 64)
}

func ConvertToFloat(raw json.RawMessage) (float64, error) {
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return 0, fmt.Errorf("failed to unmarshal raw data: %w", err)
	}
	switch v := value.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("unexpected type %T", v)
	}
}

func ResponseConvert(curRaw interface{}, curParse interface{}) error {
	curRawMap, ok := curRaw.(map[string]interface{})
	if !ok {
		return fmt.Errorf("currency data for token is not of type map[string]interface{}")
	}

	// log.Printf("resp: %+v", curRawMap)
	curRawBytes, err := json.Marshal(curRawMap)
	if err != nil {
		return fmt.Errorf("failed to marshal currency data: %w", err)
	}

	return json.Unmarshal(curRawBytes, &curParse)
}

func ConvertToWei(amount float64, decimals int) (*big.Int, error) {
	amountFloat := new(big.Float).SetFloat64(amount)

	multiplier := new(big.Float).SetFloat64(1)
	for i := 0; i < decimals; i++ {
		multiplier.Mul(multiplier, new(big.Float).SetFloat64(10))
	}

	amountWei := new(big.Float).Mul(amountFloat, multiplier)

	wei := new(big.Int)
	amountWei.Int(wei)

	return wei, nil
}

func ConvertFromWei(wei *big.Int, decimals int) float64 {
	weiFloat := new(big.Float).SetInt(wei)

	divisor := new(big.Float).SetFloat64(1)
	for i := 0; i < decimals; i++ {
		divisor.Mul(divisor, new(big.Float).SetFloat64(10))
	}

	result := new(big.Float).Quo(weiFloat, divisor)

	// Преобразуем в float64
	floatResult, _ := result.Float64()
	return floatResult
}
