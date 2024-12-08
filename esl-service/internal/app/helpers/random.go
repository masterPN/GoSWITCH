package helpers

import (
	"math/rand"
	"time"
)

// GenerateRandomFiveDigitNumber generates a random 5-digit number.
func GenerateRandomFiveDigitNumber() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(90000) + 10000 // Generates a number between 10000 and 99999
}
