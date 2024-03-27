package main

import (
	"errors"
	"fmt"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func commandHelp() error {
	fmt.Print("\nWelcome to the Pokedex!\nUsage:\n")
	structCommand := getCommandStruct()
	for _, info := range structCommand {
		fmt.Printf("\n%s: %s", info.name, info.description)
	}
	fmt.Print("\n\n")
	return nil
}

func commandExit() error {
	return errors.New("exit")
}

func parseCommand(cmd string) error {
	structCommand := getCommandStruct()
	cli, ok := structCommand[cmd]
	if !ok {
		fmt.Printf("Invalid command: %s\nUse 'help' to see the supported commands.\n", cmd)
		return nil
	} else {
		return cli.callback()
	}
}
