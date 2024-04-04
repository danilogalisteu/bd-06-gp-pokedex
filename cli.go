package main

import (
	"errors"
	"fmt"
	"internal/pokeapi"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
}

func getCommandStruct() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display next page of the list of locations",
			callback:    commandPageNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Display previous page of the list of locations",
			callback:    commandPagePrev,
		},
		"explore": {
			name:        "explore",
			description: "Show Pokemon found in the given location",
			callback:    commandExplore,
		},
	}
}

func parseCommand(cmd string) error {
	structCommand := getCommandStruct()
	command, args, _ := strings.Cut(cmd, " ")
	cli, ok := structCommand[command]
	if !ok {
		fmt.Printf("Invalid command: %s\nUse 'help' to see the supported commands.\n", cmd)
		return nil
	} else {
		return cli.callback(args)
	}
}

func commandHelp(string) error {
	fmt.Print("\nWelcome to the Pokedex!\nUsage:\n")
	structCommand := getCommandStruct()
	for _, info := range structCommand {
		fmt.Printf("\n%s: %s", info.name, info.description)
	}
	fmt.Print("\n\n")
	return nil
}

func commandExit(string) error {
	return errors.New("exiting")
}

func commandPageNext(string) error {
	locations := pokeapi.GetLocationsNext()
	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}
	return nil
}

func commandPagePrev(string) error {
	locations := pokeapi.GetLocationsPrev()
	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}
	return nil
}

func commandExplore(id string) error {
	if id == "" {
		fmt.Print("No location name was provided\n")
		return nil
	}

	encounters, ok := pokeapi.ExploreLocation(id)
	if !ok {
		fmt.Print("Given location not found in the current list:\n" + id + "\n")
		return nil
	}

	fmt.Printf("Exploring %s...\nFound Pokemon:\n", id)
	for _, encounter := range encounters.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
