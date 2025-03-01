package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Reza1878/pokedexcli/entities"
	"github.com/Reza1878/pokedexcli/internal"
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

func CommandMap(c *entities.Config, cc *internal.Cache) error {
	var response entities.LocationAreaResponse

	url := c.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	// If cache found in the cache manager, use cache otherwise make a request to POKEAPI
	if v, ok := cc.Get(url); ok {
		err := json.Unmarshal(v, &response)
		if err != nil {
			return err
		}
	} else {
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

		err = decoder.Decode(&response)
		if err != nil {
			return err
		}

		for _, r := range response.Results {
			fmt.Println(r.Name)
		}
	}

	c.Next = response.Next
	c.Previous = url

	return nil
}

func CommandMapB(c *entities.Config, cc *internal.Cache) error {
	var response entities.LocationAreaResponse

	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	// If cache found in the cache manager, use cache otherwise make a request to POKEAPI
	if v, ok := cc.Get(c.Previous); ok {
		err := json.Unmarshal(v, &response)

		if err != nil {
			return err
		}
	} else {
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

		err = decoder.Decode(&response)
		if err != nil {
			return err
		}
	}

	for _, r := range response.Results {
		fmt.Println(r.Name)
	}

	c.Next = response.Next
	c.Previous = response.Previous

	return nil
}

func CommandExplore(area string) error {
	fmt.Println("Exploring", area, "...")
	req, err := http.NewRequest("GET", fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", area), nil)
	if err != nil {
		return err
	}

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Failed to explore location")
		return errors.New("failed to explore location")
	}

	var response entities.ExploreLocationResponse
	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, val := range response.PokemonEncounters {
		fmt.Println("-", val.Pokemon.Name)
	}

	return nil
}
