package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

Repl:
	for {
		fmt.Print("pokedex > ")

		// We're scanning Stdin so no need to check if there are no tokens
		scanner.Scan()
		commandText := scanner.Text()

		switch commandText {
		case "help":
			fmt.Println("Pokedex CLI")
			fmt.Println("usage:")
			fmt.Println("	help	Show this help message.")
			fmt.Println("	exit	Exit the program.")
		case "exit":
			break Repl
		default:
			fmt.Println("Unexpected command: ", commandText)
			fmt.Println("Use the command \"help\" to see usage.")
		}

		fmt.Println()
	}
}
