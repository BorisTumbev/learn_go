package main

import (
	"os"
	"regexp"
)

func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

func SumAll(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}

func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}

	return sums
}

var digitRegexp = regexp.MustCompile("[0-9]+")

func CopyDigits(filename string) []byte {
	b, _ := os.ReadFile(filename)
	b = digitRegexp.Find(b)

	// two ways with copy and append
	// c := make([]byte, len(b))
	// copy(c, b)
	var c []byte
	c = append(c, b...)
	return c
}
