package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Reza1878/pokedexcli/entities"
	"github.com/Reza1878/pokedexcli/helper"
)

func main() {
	config := entities.Config{}

	helpCommand := entities.CliCommand{
		Name:        "help",
		Description: "Display a help message",
	}

	command := map[string]entities.CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback: func() error {
				return helper.CommandExit(&config)
			},
		},
		"help": helpCommand,
		"map": {
			Name:        "map",
			Description: "Display location area",
			Callback: func() error {
				return helper.CommandMap(&config)
			},
		},
		"mapb": {
			Name:        "mapb",
			Description: "Display previous location area",
			Callback: func() error {
				return helper.CommandMapB(&config)
			},
		},
	}

	helpCommand.Callback = func() error {
		return helper.CommandHelp(&config, command)
	}

	command["help"] = helpCommand

	commandKeys := []string{}

	for k := range command {
		commandKeys = append(commandKeys, k)
	}

	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)

		var input string

		if scanner.Scan() {
			input = scanner.Text()
		}

		cleanedInput := helper.CleanInput(input)

		for _, ck := range commandKeys {
			if ck == cleanedInput[0] {
				command[ck].Callback()
			}
		}
	}
}
