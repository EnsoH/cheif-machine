package utils

import "strconv"

func СonvertStringToFloat(line string) (float64, error) {
	return strconv.ParseFloat(line, 64)
}
