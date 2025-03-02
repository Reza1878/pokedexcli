package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Reza1878/pokedexcli/entities"
	"github.com/Reza1878/pokedexcli/helper"
	"github.com/Reza1878/pokedexcli/internal"
)

func main() {
	config := entities.Config{}
	pokedex := entities.Pokedex{
		Pokemon: make(map[string]entities.Pokemon),
		Tries:   make(map[string]int, 0),
	}
	cacheManager := internal.NewCache(5 * time.Second)

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
				return helper.CommandMap(&config, cacheManager)
			},
		},
		"mapb": {
			Name:        "mapb",
			Description: "Display previous location area",
			Callback: func() error {
				return helper.CommandMapB(&config, cacheManager)
			},
		},
		"explore": {
			Name:        "explore",
			Description: "Explore specific location",
		},
		"catch": {
			Name:        "catch",
			Description: "Catch a Pokemon",
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

		if len(cleanedInput) > 0 {
			if cleanedInput[0] == "explore" && len(cleanedInput) > 1 {
				err := helper.CommandExplore(cleanedInput[1])
				if err != nil {
					fmt.Println(err)
				}

			} else if cleanedInput[0] == "catch" && len(cleanedInput) > 1 {
				err := helper.CommandCatch(&pokedex, cleanedInput[1])
				if err != nil {
					fmt.Println(err)
				}

			} else if cleanedInput[0] == "inspect" && len(cleanedInput) > 1 {
				err := helper.CommandInspect(&pokedex, cleanedInput[1])
				if err != nil {
					fmt.Println(err)
				}

			} else {
				for _, ck := range commandKeys {
					if ck == cleanedInput[0] {
						command[ck].Callback()
					}
				}
			}
		}
	}
}
