package main

import (
	"fmt"
	"strings"
)

// Empty struct
type Unit struct{}

// This is exported (public) because it starts with a capital letter
type Point struct {
	x int
	y int
}

func printPoint(p Point) {
	fmt.Printf("Point{%v, %v}\n", p.x, p.y)
}

func pointers() {
	// var (
	// 	x int = 5
	// 	p *int = &x
	// )

	// var x int = 5
	// var p *int

	// Pointers
	x := 5
	p := &x

	// Casting to pointer needs parentheses around the type
	y := (*int)(nil)

	fmt.Printf("p: %v\n", p)
	fmt.Printf("p == y: %v\n", p != y)

	fmt.Printf("*p: %v\n", *p)

	// No pointer arithmetic, so this is just (*p)++
	*p++
	fmt.Printf("x: %v\n", x)
	fmt.Println()
}

func structs() {
	// Prints {}
	fmt.Println(Unit{})

	myPoint := Point{3, 4}
	printPoint(myPoint)

	// Pointer to a struct
	p := &myPoint

	// Access struct fields by dereferencing
	(*p).x = 6

	// Or without explicitly dereferencing (no need for arrow operator)
	p.y *= 2

	printPoint(myPoint)

	// Struct initialization (literals)
	var (
		v1 = Point{1, 2}  // has type Point
		v2 = Point{x: 1}  // y:0 is implicit, order is not important
		v3 = Point{}      // x:0 and y:0
		p1 = &Point{1, 2} // Directly take pointer to the value
	)

	fmt.Println(v1, p1, v2, v3)

	// This is an anonymous (unnamed) struct
	anon := struct {
		x int
		y int
	}{1, 2}

	fmt.Printf("anon is %v of type %T\n", anon, anon)

	fmt.Println()
}

func arrays() {
	var arr1 [10]int
	arr1[0] = 42
	arr1[len(arr1)-1] = 67
	fmt.Println(arr1[0])

	// This won't work because arr has type [10]int (fixed size)
	// arr = [2]int{}

	// Here, you can put a larger size than the initializer list
	// and the rest of the elements will have zero values
	arr2 := [3]int{1, 2, 3}
	fmt.Println(arr2)

	// This slice will have indices 1 and 2 from arr2
	var slice2 []int = arr2[1:3]
	fmt.Println(slice2)

	// Slice of the entire array (low is 0 and high is len(arr1))
	slice1 := arr1[:]
	fmt.Println(slice1)

	// You can even reassign with a different size
	slice1 = arr1[7:]
	fmt.Println(slice1)

	// Slices are just references to arrays
	// This accesses arr2[1] since slice2 is arr2[1:3]
	slice2[0] = 80
	fmt.Println(arr2)

	// This creates an array then builds a slice that references it
	slice3 := []bool{true, false, true}
	fmt.Println(slice3[1:])

	fmt.Println()
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func sliceCapacity() {
	// Length is len of slice, capacity is length of underlying array of slice, starting from first element of slice
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	// Slice the slice to give it zero length.
	s = s[:0]
	printSlice(s)

	// Extend its length.
	s = s[:4]
	printSlice(s)

	// Drop its first two values.
	// This will change the capacity. Earlier we were able to change len without cap
	// because we did s[:0] (limit high bound) instead of s[2:] (limit low bound)
	s = s[2:]
	printSlice(s)

	// This will not change the low bound at all (obviously). We don't have negative indices here like in Python
	s = s[0:]
	printSlice(s)

	// nil slice, length and capacity are 0
	var s2 []int
	fmt.Println(s2, len(s2), cap(s2))
	fmt.Printf("s2 == nil: %v\n", s2 == nil)

	// Dynamically sized arrays: use make
	s3 := make([]int, len(s)+2)
	fmt.Println(s3)

	// Can set higher capacity
	s4 := make([]int, len(s)+2, len(s)+2+10)
	fmt.Println(s4)

	// Extend s4 to it's max capacity
	s4 = s4[:cap(s4)]
	fmt.Println(s4)

	// Range iterator
	for index, value := range s[:cap(s)] {
		fmt.Printf("index %v: %v\n", index, value)
	}

	// New to Go 1.22
	for i := range 5 {
		fmt.Println(i * i)
	}

	fmt.Println()
}

// Check https://go.dev/tour/moretypes/18
func Pic(dx, dy int) [][]uint8 {
	arr := make([][]uint8, dx)
	for i := range dx {
		arr[i] = make([]uint8, dy)
		for j := range dy {
			arr[i][j] = uint8(i ^ j)
		}
	}

	return arr
}

func maps() {
	// Map from string to int
	var m1 map[string]int

	// Need to allocate
	m1 = make(map[string]int)

	// Insert elements
	m1["a"] = 4
	m1["c"] = 3

	fmt.Println(m1)

	// Initializer (map literal)
	m2 := map[string]Point{
		"hello": {1, 2},
		"world": {3, 4}, // Type name can be omitted
		"!":     {5, 6}, // terminating comma is required
	}

	fmt.Println(m2)

	// Remove by key (this won't panic or return an error if the key doesn't exist)
	delete(m2, "!")

	// Retrieve element
	val := m2["hello"]

	// Check if key exists
	val, ok := m2["hello"]
	if ok {
		fmt.Println(val)
	}

	val, ok = m2["!"]
	if !ok {
		// In this case val will have zero value
		fmt.Println("! not found in m2")
	}

	fmt.Println(m2)

	fmt.Println()
}

// Check https://go.dev/tour/moretypes/23
func WordCount(s string) map[string]int {
	wordMap := make(map[string]int)

	// Fields is basically str.split in Python
	// Iterate over words after splitting by whitespace
	for _, word := range strings.Fields(s) {
		// If key does not exist, insert it
		if _, ok := wordMap[word]; !ok {
			wordMap[word] = 0
		}

		wordMap[word]++
	}

	return wordMap
}

// Functions can be passed as parameters
func foo(bar func(x int) bool) string {
	if bar(42) {
		return "foobar"
	}

	return "foobaz"
}

func adder() func(int) int {
	// This variable here will be "bound" to the returned function
	// We basically return a function with *state*, without having to use global variables or parameters
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func fibonacci() func() int {
	// We can make use of closures to store the fibonacci state
	f0 := 0
	f1 := 0

	return func() int {
		if f0 == f1 && f0 == 0 {
			// Base case 1 (return 0)
			f0 = 1
			return 0
		} else if f1 == 0 && f0 == 01 {
			// Base case 2 (return 1)
			f0 = 0
			f1 = 1
			return 1
		}

		// New fib
		f2 := f1 + f0

		// Update old fibs
		f0 = f1
		f1 = f2

		return f2
	}
}

func functions() {
	myBar := func(x int) bool {
		return x > 40
	}

	fmt.Println(foo(myBar))

	x := 5

	// A closure is a function value that references variables outside its body
	baz := func() {
		fmt.Println(x)
	}

	baz() // Prints 5
	x++   // Update x
	baz() // Now prints 6

	fmt.Println("\n---- Adder example")

	// Each adder will act independently
	// Each of them has its own sum variable (basically internal state)
	pos, neg := adder(), adder()
	for i := 0; i < 4; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}

	fmt.Println("\n---- Fibonacci example")

	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

func main() {
	fmt.Println("---- Pointers ----")
	pointers()

	fmt.Println("---- Structs ----")
	structs()

	fmt.Println("---- Arrays and Slices ----")
	arrays()

	fmt.Println("---- Slice capacity ----")
	sliceCapacity()

	fmt.Println("---- Maps ----")
	maps()

	fmt.Println("---- Functions ----")
	functions()
}
