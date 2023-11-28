package main

import (
	"fmt"
	"math/rand"
)

func moveDecimalRight(amount float64) string {
	return fmt.Sprint(int(amount * 100))
}

func padLeftWithZeros(input string, length int) string {
	return fmt.Sprintf("%0*s", length, input)
}

func generateRrn() string {
	min := 0
	max := 999999999999
	randomTwelveDigitNumber := rand.Intn(max-min+1) + min
	return fmt.Sprintf("%012d", randomTwelveDigitNumber)
}

func addTrailingSpaces(input string, length int) string {
	return fmt.Sprintf("%-*s", length, input)
}
