package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

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

func catchProbability(pokemon entities.Pokemon) float64 {
	return 1 / (1 + float64(pokemon.BaseExperience)/100)
}

func tryCatch(p *entities.Pokedex, pokemon entities.Pokemon) bool {
	rand.Seed(time.Now().UnixNano())
	probability := catchProbability(pokemon)

	temp := rand.Float64()

	temp -= float64(p.Tries[pokemon.Name]) * float64(0.015)

	return temp < probability
}

func CommandCatch(p *entities.Pokedex, pokemonName string) error {
	fmt.Printf("Throwing a Pokeball at %s\n", pokemonName)
	res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName))
	if err != nil {
		return fmt.Errorf("failed to get pokemon info: %v", err)
	}
	defer res.Body.Close()

	var pokemon entities.Pokemon
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pokemon)

	if err != nil {
		return fmt.Errorf("failed to decode json data: %v", err)
	}

	if _, ok := p.Tries[pokemonName]; ok {
		p.Tries[pokemonName] += 1
	} else {
		p.Tries[pokemonName] = 0
	}

	isCatched := tryCatch(p, pokemon)

	if isCatched {
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
		p.Pokemon[pokemonName] = pokemon
		p.Tries[pokemonName] = 0
	} else {
		fmt.Printf("Failed to catch %s\n", pokemonName)
	}

	return nil
}

func CommandInspect(p *entities.Pokedex, pokemonName string) error {
	if _, ok := p.Pokemon[pokemonName]; !ok {
		return errors.New("you have not cought that pokemon")
	}

	currPokemon := p.Pokemon[pokemonName]

	fmt.Printf("Name: %s\n", pokemonName)
	fmt.Printf("Height: %v\n", currPokemon.Height)
	fmt.Printf("Weight: %v\n", currPokemon.Weight)
	fmt.Println("Stats:")

	for _, s := range currPokemon.Stats {
		fmt.Printf("  -%s: %v\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range currPokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}
