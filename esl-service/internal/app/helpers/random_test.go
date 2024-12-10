package helpers

import (
	"testing"
)

func TestGenerateRandomFiveDigitNumber(t *testing.T) {
	// Test that the generated number is within the expected range
	num := GenerateRandomFiveDigitNumber()
	if num < 10000 || num > 99999 {
		t.Errorf("Generated number %d is out of range", num)
	}

	// Test that the generated number is different each time the function is called
	num1 := GenerateRandomFiveDigitNumber()
	num2 := GenerateRandomFiveDigitNumber()
	if num1 == num2 {
		t.Errorf("Generated numbers %d and %d are the same", num1, num2)
	}

	// Test that the generated number is not less than 10000
	num = GenerateRandomFiveDigitNumber()
	if num < 10000 {
		t.Errorf("Generated number %d is less than 10000", num)
	}

	// Test that the generated number is not greater than 99999
	num = GenerateRandomFiveDigitNumber()
	if num > 99999 {
		t.Errorf("Generated number %d is greater than 99999", num)
	}
}

func TestGenerateRandomFiveDigitNumberConsistency(t *testing.T) {
	// Test that the generated numbers are consistent across multiple calls
	nums := make([]int, 100)
	for i := range nums {
		nums[i] = GenerateRandomFiveDigitNumber()
	}
	for _, num := range nums {
		if num < 10000 || num > 99999 {
			t.Errorf("Generated number %d is out of range", num)
		}
	}
}
