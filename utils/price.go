package utils

import (
	"fmt"
	"strings"
)

func FormatPrice(price float64) string {
	s := fmt.Sprintf("%.2f", price)

	parts := strings.Split(s, ".")
	intPart := parts[0]
	decimalPart := parts[1]

	var result string
	for i := len(intPart); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}

		if len(result) > 0 {
			result = intPart[start:i] + " " + result
		} else {
			result = intPart[start:i]
		}
	}

	return result + "," + decimalPart
}
