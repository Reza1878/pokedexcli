package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Reza1878/pokedexcli/entities"
)

func CleanInput(text string) []string {
	words := []string{}

	for _, v := range strings.Split(text, " ") {
		if v != "" {
			words = append(words, strings.ToLower(strings.Trim(v, " ")))
		}
	}

	return words
}

func CommandExit(c *entities.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	os.Exit(0)

	return nil
}

func CommandHelp(c *entities.Config, command map[string]entities.CliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for _, v := range command {
		fmt.Printf("%s: %s\n", v.Name, v.Description)
	}

	return nil
}

func CommandMap(c *entities.Config) error {
	url := c.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var response entities.LocationAreaResponse
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	for _, r := range response.Results {
		fmt.Println(r.Name)
	}

	c.Next = response.Next
	c.Previous = url

	return nil
}

func CommandMapB(c *entities.Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	req, err := http.NewRequest("GET", c.Previous, nil)
	if err != nil {
		return nil
	}

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var response entities.LocationAreaResponse
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	for _, r := range response.Results {
		fmt.Println(r.Name)
	}

	c.Next = response.Next
	c.Previous = response.Previous

	return nil
}
