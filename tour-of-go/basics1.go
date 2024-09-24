package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"runtime"
	"strconv"
)

// Multiple var declarations
var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

// Function returning two values (also naked)
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x

	return
}

func hello(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// Naked return, aka "copy out"
func nakedReturn() (x int) {
	x = 5

	return
}

func Sqrt(x float64) float64 {
	var z float64 = x / 2
	zprev := float64(0)

	// Newton's method
	for math.Abs(zprev-z) > 0.00000001 {
		zprev = z
		z -= (z*z - x) / (2 * z)
	}

	return z
}

func main1() {
	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)

	// Initialization
	sum := func(a, b int) int { return a + b }(3, 4)

	// Const needs a value (type is optional)
	const x float32 = 12.24

	// Declaration of a function
	//
	// Go's type declaration fixes issues you didn't know you had in C:
	// https://go.dev/blog/declaration-syntax
	var isFloatGreater func(float32, int) bool

	// Assign a function (closure) to the variable
	// Notice that the only difference is that the function name is omitted
	isFloatGreater = func(f float32, i int) bool {
		return f > float32(i)
	}

	fmt.Printf("Sum: %v\n", sum)
	fmt.Printf("pow: %v\n", isFloatGreater(x, 2))
	fmt.Printf("Pow: %v\n", math.Pow(float64(x), 2))

	hello("world")

	var i int = 5

	// Normal for loop
	for i = 0; i < 10; i++ {
	}

	// No init statement
	for ; i > 0; i-- {
	}

	// No init or post statements
	for i == 42 {
	}

	i = 4

	// while(true)
	for {

		// Init statement for if block
		if i++; i < 4 {
			println("i: " + strconv.Itoa(i) + " is less than 4")
		} else if i++; i > 6 {
			println("i: " + strconv.FormatInt(int64(i), 10) + " is greater than 6")
		} else if x := true; i == 0 {
			// x is declared and in scope only for this "else if" block and the blocks after (the else)
			break
		} else {
			println(x)
			break
		}

		println("early return")
		return
	}

	// Note: println prints to stderr, so any calls before here will actually be printed after this (and later) printf
	fmt.Printf("Sqrt: %v\n", float64(Sqrt(64)))

	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}
}
