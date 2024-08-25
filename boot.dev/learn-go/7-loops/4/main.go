package main

import "fmt"

func fizzbuzz() {
	for i := 1; i <= 100; i++ {
		divBy3 := i%3 == 0
		divBy5 := i%5 == 0

		if divBy3 && divBy5 {
			fmt.Println("fizzbuzz")
		} else if divBy3 {
			fmt.Println("fizz")
		} else if divBy5 {
			fmt.Println("buzz")
		} else {
			fmt.Println(i)
		}
	}
}

// don't touch below this line

func main() {
	fizzbuzz()
}
