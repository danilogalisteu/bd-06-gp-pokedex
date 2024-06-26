package main

import (
	"errors"
	"fmt"
	"internal/pokeapi"
	"strings"
	"math/rand"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
}

var pokedex = make(map[string]pokeapi.PokeInfo)

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
			description: "Show all Pokemon found in the given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch the given Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "See details of the given Pokemon, if it was caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See a list of the Pokemon caught",
			callback:    commandPokedex,
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
		fmt.Print("Given location not found in the current list: " + id + "\n")
		return nil
	}

	fmt.Printf("Exploring %s...\nFound Pokemon:\n", id)
	for _, encounter := range encounters.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(id string) error {
	if id == "" {
		fmt.Print("No Pokemon name was provided\n")
		return nil
	}

	info, ok := pokeapi.CatchPokemon(id)
	if !ok {
		fmt.Print("Given Pokemon not found in the current location: " + id + "\n")
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", id)

	if rand.Intn(1000) < info.BaseExperience {
		fmt.Printf("%s escaped!\n", id)
		return nil
	}

	fmt.Printf("%s was caught!\nYou may now inspect it with the inspect command.\n", id)
	pokedex[id] = info

	return nil
}

func commandInspect(id string) error {
	if id == "" {
		fmt.Print("No Pokemon name was provided\n")
		return nil
	}

	info, ok := pokedex[id]
	if !ok {
		fmt.Printf("Given Pokemon is not in the Pokedex: %s\n", id)
		return nil
	}

	fmt.Printf("Name: %s\n", id)
	fmt.Printf("Height: %d\n", info.Height)
	fmt.Printf("Weight: %d\n", info.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range info.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, tp := range info.Types {
		fmt.Printf("  -%s\n", tp.Type.Name)
	}

	return nil
}

func commandPokedex(string) error {
	fmt.Print("Your Pokedex:\n")
	for _, info := range pokedex {
		fmt.Printf(" - %s\n", info.Name)
	}
	return nil
}
