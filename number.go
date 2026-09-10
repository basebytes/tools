package tools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func NumberToShortPrice(number int) (price string) {
	price = strings.TrimRight(fmt.Sprintf("%.2f", float64(number)/100), "0")
	price = strings.TrimRight(price, ".")
	return
}

func NumberToFixedPrice(number int) string {
	return fmt.Sprintf("%.2f", float64(number)/100)
}

func PriceToNumber(price string) int {
	c, cents := 0, -1
	prices := strings.Split(price, ".")
	y, err := strconv.Atoi(prices[0])
	if err == nil && len(prices) > 1 {
		if len(prices[1]) > 2 {
			c, err = strconv.Atoi(prices[1][:2])
		} else {
			c, err = strconv.Atoi(prices[1])
		}
	}
	if err == nil {
		cents = y*100 + c
	}
	return cents
}

func ToNumber(numStr string) int {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		num = -1
	}
	return num
}

func Trunc(value float64, precision int) float64 {
	multiples := math.Pow10(precision)
	return math.Trunc(value*multiples) / multiples
}
