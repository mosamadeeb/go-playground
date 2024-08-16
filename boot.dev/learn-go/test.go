package main

import "fmt"

type Vertex struct {
	X int
	Y int
}

func main() {
	const name = "Saul Goodman"
	const openRate = 30.5
	var test = 13
	fmt.Println(test)

	msg := fmt.Sprintf("Hi %+v, your open rate is %#v percent\n", Vertex{1, 2}, openRate)

	// don't edit below this line

	fmt.Print(msg)
}
