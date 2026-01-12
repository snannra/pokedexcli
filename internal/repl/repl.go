package repl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/snannra/pokedexcli/internal/pokecache"
)

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config, *string) error
}

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func CommandExit(cfg *Config, _ *string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

type locationAreaListResp struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

type pokemonListResp struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func CommandMap(cfg *Config, _ *string) error {
	const baseLocationAreaURL = "https://pokeapi.co/api/v2/location-area/"

	url := cfg.Next
	if url == "" {
		url = baseLocationAreaURL
	}

	if entry, exists := cfg.Cache.Get(url); exists {
		var out locationAreaListResp
		if err := json.Unmarshal(entry, &out); err != nil {
			return fmt.Errorf("unmarshal cached response: %w", err)
		}
		for _, r := range out.Results {
			fmt.Println(r.Name)
		}
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get location areas: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read respones: %w", err)
	}

	var out locationAreaListResp
	if err := json.Unmarshal(data, &out); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	for _, r := range out.Results {
		fmt.Println(r.Name)
	}

	cfg.Cache.Add(url, data)

	if out.Next != nil {
		cfg.Next = *out.Next
	} else {
		cfg.Next = ""
	}

	if out.Previous != nil {
		cfg.Previous = *out.Previous
	} else {
		cfg.Previous = ""
	}

	return nil
}

func CommandMapB(cfg *Config, _ *string) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	url := cfg.Previous

	if entry, exists := cfg.Cache.Get(url); exists {
		var out locationAreaListResp
		if err := json.Unmarshal(entry, &out); err != nil {
			return fmt.Errorf("unmarshal cached response: %w", err)
		}
		for _, r := range out.Results {
			fmt.Println(r.Name)
		}
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get location areas: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read respones: %w", err)
	}

	var out locationAreaListResp
	if err := json.Unmarshal(data, &out); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	for _, r := range out.Results {
		fmt.Println(r.Name)
	}

	cfg.Cache.Add(url, data)

	if out.Next != nil {
		cfg.Next = *out.Next
	} else {
		cfg.Next = ""
	}

	if out.Previous != nil {
		cfg.Previous = *out.Previous
	} else {
		cfg.Previous = ""
	}

	return nil
}

func CommandExplore(cfg *Config, location *string) error {
	baseLocationAreaURL := "https://pokeapi.co/api/v2/location-area/" + *location + "/"

	fmt.Printf("Exploring %s...\nFound Pokemon:\n", *location)
	if entry, exists := cfg.Cache.Get(baseLocationAreaURL); exists {
		var out pokemonListResp
		if err := json.Unmarshal(entry, &out); err != nil {
			return fmt.Errorf("unmarshal cached response: %w", err)
		}

		for _, p := range out.PokemonEncounters {
			fmt.Println(" - " + p.Pokemon.Name)
		}

		return nil
	}

	resp, err := http.Get(baseLocationAreaURL)
	if err != nil {
		return fmt.Errorf("get pokemon encounters: %w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	cfg.Cache.Add(baseLocationAreaURL, data)

	var out pokemonListResp
	if err := json.Unmarshal(data, &out); err != nil {
		return fmt.Errorf("unmarshal cached response: %w", err)
	}

	for _, p := range out.PokemonEncounters {
		fmt.Println(" - " + p.Pokemon.Name)
	}

	return nil
}
