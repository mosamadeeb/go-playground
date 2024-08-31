package main

import (
	"bufio"
	"fmt"
	"os"
)

type command struct {
	name     string
	usage    string
	callback func() error
}

// This is used to determine the order of enumeration of the commands
var commandList = []string{
	"help", "exit",
}

// An initializer is not used to avoid a circular dependency with the [commandHelp] function
var commandMap map[string]command = make(map[string]command)

func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, com := range commandList {
		fmt.Printf("%s: %s\n", com, commandMap[com].usage)
	}

	return nil
}

func commandExit() error {
	// Technically, a defer here makes no sense because the return value won't reach the caller anyway
	defer os.Exit(0)
	return nil
}

func initCommandMap() {
	commandMap["help"] = command{"help", "Displays a help message", commandHelp}
	commandMap["exit"] = command{"exit", "Exit the Pokedex", commandExit}
}

func main() {
	initCommandMap()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// We're scanning Stdin so no need to check if there are no tokens
		scanner.Scan()
		commandText := scanner.Text()

		com, ok := commandMap[commandText]
		if ok {
			com.callback()
		} else {
			fmt.Println("Unexpected command: ", commandText)
			fmt.Println("Use the command \"help\" to see usage")
		}

		fmt.Println()
	}
}
