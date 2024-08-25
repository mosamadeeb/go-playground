package main

import (
	"fmt"
)

func printPrimes(max int) {
	// Print the only even prime number
	fmt.Println(2)

	// Go over all odd integers up until and including max
	for n := 3; n <= max; n += 2 {
		// n is a prime number until proven otherwise
		isPrime := true

		for i := 3; i < n; i++ {
			if n%i == 0 {
				// n is not a prime number
				isPrime = false
				break
			}
		}

		if isPrime {
			fmt.Println(n)
		}
	}
}

// don't edit below this line

func test(max int) {
	fmt.Printf("Primes up to %v:\n", max)
	printPrimes(max)
	fmt.Println("===============================================================")
}

func main() {
	test(10)
	test(20)
	test(30)
}
