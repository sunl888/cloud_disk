package tests

import "testing"

func TestSum(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	expected := 15
	actual := Sum(numbers)

	if actual != expected {
		t.Errorf("Expected the sum of %v to be %d but instead got %d!", numbers, expected, actual)

	}
}

func Sum(numbers []int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}
