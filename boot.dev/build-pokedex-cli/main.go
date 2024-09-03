package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mosamadeeb/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name     string
	usage    string
	callback func(args []string) error
}

// An initializer is not used to avoid a circular dependency with the [commandHelp] function
var cliCommandMap map[string]cliCommand = make(map[string]cliCommand)

var mapPageConfig pokeapi.PageConfig

func commandMap(args []string, prevPage bool) error {
	if len(args) != 0 {
		return errors.New("command takes no arguments")
	}

	areas, err := pokeapi.FetchLocationArea(&mapPageConfig, prevPage)
	if err != nil {
		return fmt.Errorf("could not fetch locations: %w", err)
	}

	for _, s := range areas {
		fmt.Println(s)
	}

	return nil
}

func commandExplore(args []string) error {
	if len(args) != 1 {
		return errors.New("command takes only 1 argument")
	}

	locationArea := args[0]
	fmt.Printf("Exploring %s...\n", locationArea)

	pokemon, err := pokeapi.QueryLocationAreaPokemon(locationArea)
	if err != nil {
		return fmt.Errorf("could not query pokemon: %w", err)
	}

	fmt.Println("Found Pokemon:")
	for _, p := range pokemon {
		fmt.Println(" -", p)
	}

	return nil
}

func commandCatch(args []string) error {
	if len(args) != 1 {
		return errors.New("command takes only 1 argument")
	}

	name := args[0]
	pokemon, err := pokeapi.QueryPokemon(name)
	if err != nil {
		return fmt.Errorf("could not query pokemon: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	if pokemon.TryCatch() {
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandHelp(args []string) error {
	if len(args) != 0 {
		return errors.New("command takes no arguments")
	}

	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, com := range cliCommandList {
		fmt.Printf("%s: %s\n", cliCommandMap[com].name, cliCommandMap[com].usage)
	}

	return nil
}

func commandExit(args []string) error {
	if len(args) != 0 {
		return errors.New("command takes no arguments")
	}

	// Technically, a defer here makes no sense because the return value won't reach the caller anyway
	defer os.Exit(0)
	return nil
}

// This is used to determine the order of enumeration of the commands
var cliCommandList = []string{"map", "mapb", "explore", "catch", "help", "exit"}

// This built-in feature allows us to do things before the main function is executed
// This is run only once per package, but we can have as many init() functions as we want and they will all be executed
// Pretty cool for simple cases like this
func init() {
	cliCommandMap["map"] = cliCommand{"map", "Displays the next 20 location areas", func(args []string) error {
		return commandMap(args, false)
	}}
	cliCommandMap["mapb"] = cliCommand{"mapb", "Displays the previous 20 location areas", func(args []string) error {
		return commandMap(args, true)
	}}
	cliCommandMap["explore"] = cliCommand{"explore <area_name>", "Displays Pokemon found in a location area", commandExplore}
	cliCommandMap["catch"] = cliCommand{"catch <pokemon_name>", "Attempts to catch a Pokemon", commandCatch}
	cliCommandMap["help"] = cliCommand{"help", "Displays a help message", commandHelp}
	cliCommandMap["exit"] = cliCommand{"exit", "Exit the Pokedex", commandExit}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// We're scanning Stdin so no need to check if there are no tokens
		scanner.Scan()

		scannedText := strings.Split(scanner.Text(), " ")
		commandText := scannedText[0]
		commandArgs := scannedText[1:]

		com, ok := cliCommandMap[commandText]
		if ok {
			err := com.callback(commandArgs)
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