package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

// Method defined on Vertex
// The (v Vertex) part is called a receiver argument
// This is technically a const method as it can't modify the struct v
func (v Vertex) Mirrored() Vertex {
	return Vertex{v.Y, v.X}
}

// This is method with a pointer receiver
// Basically a non-const method
func (v *Vertex) Translate(x, y float64) {
	v.X += x
	v.Y += y
}

// Basically a typedef
type MyFloat float64

// We can't define a method on float64 directly because it wasn't
// defined in this package. So we use the new type MyFloat instead
func (f MyFloat) Squared() MyFloat {
	return f * f
}

func methods() {
	var v Vertex

	v = Vertex{2, 4}
	fmt.Println(v.Mirrored())

	// Translate has a pointer receiver, but we can call it directly on v
	// Interpreted as (&v).Translate()
	v.Translate(6, 2)
	fmt.Println(v)

	// Similarly, Mirrored is a value receiver, but we can call it on a pointer
	// Interpreted as (*p).Mirrored()
	p := &v
	fmt.Println(p.Mirrored())

	fmt.Println(MyFloat(2.4).Squared())
	fmt.Println()
}

// An interface defining two methods
type geometric interface {
	Perimeter() float64
	Area() float64
}

// We cannot have an interface as a receiver for a method
// func (g geometric) print() {}

// So how do we implement a method that is defined on an interface, with the same
// functionality for all implementations of that interface?

// We should instead write a function that takes the interface as an argument
// It can take a value or pointer
func print(g geometric) {
	fmt.Println("Geometric\n--------------")
	fmt.Println(g)
	fmt.Printf("Perimeter: %v\n", g.Perimeter())
	fmt.Printf("Area: %v\n\n", g.Area())
}

type square struct {
	length float64
}

// Perimeter is implemented for *square
func (s *square) Perimeter() float64 {
	if s == nil {
		return 0
	}

	return s.length * 4
}

// Area is implemented for *square
func (s *square) Area() float64 {
	if s == nil {
		return 0
	}

	return s.length * s.length
}

// This means that *square implements the interface.
// There is no "explicit intent". Implicit implementation of interfaces
// decouples the definition of the interface from its implementation
// This means that package A with an interface does not have to be arranged
// before package B with a type that implements that interface

// Another type that implements this interface
type circle struct {
	radius float64
}

// Methods are implemented by value instead of pointer
func (c circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func (c circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

// This is a method for the Stringer interface
func (c circle) String() string {
	// Be careful here because if you format c as %v or do println then this function will be called again
	// which would lead to an infinite loop
	return fmt.Sprintf("circle with radius = %v", c.radius)
}

// Function that accepts an empty interface
func printlnWhatever(x interface{}) {
	if x != nil {
		fmt.Printf("(%v, %T)\n", x, x)
	} else {
		fmt.Println("x is nil")
	}
}

func interfaces() {
	var g geometric

	// This won't work because square (value) does not implement the interface
	// g = square{5}

	// This will work though
	g = &square{5}
	print(g)

	g = circle{radius: 7}
	print(g)

	// If we set g to the pointer of a circle, it would also work
	g = &circle{radius: 7}
	print(g)

	// Let's try messing with nil values
	// Interface is defined on *square
	var s *square

	// This works, and it just passes the pointer to the method
	s = nil
	fmt.Println(s.Perimeter())

	// This also works
	g = s
	fmt.Println(g.Perimeter())

	// We can't do the same with circle because its methods only pass by value
	// This means we have two types that implement the same interface,
	// But one of them can cause a panic at runtime because of a nil dereference

	// But honestly, a value receiver in this case is just like a pointer receiver which we don't check for nil

	// We can think of interface values as a tuple of value and concrete type
	// (value, type)
	//
	// If the *concrete value* is nil (like in the above cases), the type can still be used
	// So the method will be called with a nil receiver
	// The *interface value* is not actually nil in this case
	//
	// If the interface value is nil, this means that it neither has a value nor a type
	// This means we can't even call any methods
	//g = nil
	//g.Perimeter()
	// ^ The above will panic

	fmt.Println()

	// This variable is an empty interface
	// which means it can contain any value that implements at least zero methods (basically anything)
	// This can be used to handle values of unknown types
	var e1 interface{}
	printlnWhatever(e1)

	e1 = 42
	printlnWhatever(e1)

	e1 = 3.14
	printlnWhatever(e1)

	e1 = &circle{7}
	printlnWhatever(e1)

	fmt.Println()

	// We can assert the type of an interface value
	// If the type is different, we will get a runtime error like this
	// panic: interface conversion: interface {} is *main.circle, not X
	e1Circle := e1.(*circle)
	fmt.Printf("e1Circle: %v\n", e1Circle)

	// If we add a second variable on the left hand side to store the success of the assertion,
	// then the assertion will not cause a panic if it fails
	e1Int, ok := e1.(int)

	if ok {
		fmt.Printf("e1Int: %v", e1Int)
	} else {
		fmt.Println("e1 is not an int!")
	}

	// We can also assert types without casting the value if we want to
	// Again, this will panic if the type is wrong
	_ = e1.(*circle)

	// or we can use the test boolean
	if _, ok = e1.(float32); !ok {
		fmt.Println("e1 is not an float32!")
	}

	// We can also switch on the type of an interface value by using the keyword type
	switch v := e1.(type) {
	case circle:
		//
	case *circle:
		// It should be this
		fmt.Println(v)
	case int:
		// Remember, case blocks in Go do not fallthrough
		// unless we use the keyword *fallthrough*, although it can't be used when switching on a type
	default:
		// If we reach this, then v will have the same type as e1
	}

	// circle has the String method of the interface Stringer implemented
	// Now anything that formats a circle as %v will use the method instead
	fmt.Println(e1.(*circle))

	fmt.Println()
}

func main3() {
	fmt.Println("---- Methods ----")
	methods()

	fmt.Println("---- Interfaces ----")
	interfaces()
}
