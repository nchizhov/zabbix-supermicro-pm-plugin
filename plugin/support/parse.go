package support

import "strconv"

func GetDigit(digitData []byte) (float64, error) {
	digit, err := strconv.ParseFloat(string(digitData), 64)
	if err == nil {
		return digit, nil
	}
	return 0, err
}
