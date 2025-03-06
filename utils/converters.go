package utils

import (
	"encoding/json"
	"fmt"
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

	curRawBytes, err := json.Marshal(curRawMap)
	if err != nil {
		return fmt.Errorf("failed to marshal currency data: %w", err)
	}

	return json.Unmarshal(curRawBytes, &curParse)
}
