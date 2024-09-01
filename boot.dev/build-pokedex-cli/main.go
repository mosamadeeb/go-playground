package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mosamadeeb/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name     string
	usage    string
	callback func() error
}

// An initializer is not used to avoid a circular dependency with the [commandHelp] function
var cliCommandMap map[string]cliCommand = make(map[string]cliCommand)

var mapPageConfig pokeapi.PageConfig

func commandMap(prevPage bool) error {
	areas, err := pokeapi.FetchLocationArea(&mapPageConfig, prevPage)
	if err != nil {
		return fmt.Errorf("could not fetch locations: %w", err)
	}

	for _, s := range areas {
		fmt.Println(s)
	}

	return nil
}

func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, com := range cliCommandList {
		fmt.Printf("%s: %s\n", com, cliCommandMap[com].usage)
	}

	return nil
}

func commandExit() error {
	// Technically, a defer here makes no sense because the return value won't reach the caller anyway
	defer os.Exit(0)
	return nil
}

// This is used to determine the order of enumeration of the commands
var cliCommandList = []string{"map", "mapb", "help", "exit"}

// This built-in feature allows us to do things before the main function is executed
// This is run only once per package, but we can have as many init() functions as we want and they will all be executed
// Pretty cool for simple cases like this
func init() {
	cliCommandMap["map"] = cliCommand{"map", "Displays the next 20 location areas", func() error {
		return commandMap(false)
	}}
	cliCommandMap["mapb"] = cliCommand{"mapb", "Displays the previous 20 location areas", func() error {
		return commandMap(true)
	}}
	cliCommandMap["help"] = cliCommand{"help", "Displays a help message", commandHelp}
	cliCommandMap["exit"] = cliCommand{"exit", "Exit the Pokedex", commandExit}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// We're scanning Stdin so no need to check if there are no tokens
		scanner.Scan()
		commandText := scanner.Text()

		com, ok := cliCommandMap[commandText]
		if ok {
			err := com.callback()
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unexpected command: ", commandText)
			fmt.Println("Use the command \"help\" to see usage")
		}

		fmt.Println()
	}
}
