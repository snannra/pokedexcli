package repl

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/snannra/pokedexcli/internal/pokecache"
)

type Pokemon struct {
	BaseExperience int `json:"base_experience"`
	Height         int `json:"height"`
	Weight         int `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	PokeDex  map[string]Pokemon
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

func CommandCatch(cfg *Config, pokemonName *string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", *pokemonName)

	pokemonURL := "https://pokeapi.co/api/v2/pokemon/" + *pokemonName + "/"

	req, err := http.Get(pokemonURL)
	if err != nil {
		return fmt.Errorf("get pokemon data: %w", err)
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("read pokemon data: %w", err)
	}

	var pokeData Pokemon
	if err := json.Unmarshal(data, &pokeData); err != nil {
		return fmt.Errorf("unmarshal pokemon data: %w", err)
	}

	threshold := pokeData.BaseExperience - 70
	catchChance := rand.Intn(pokeData.BaseExperience)

	if catchChance > threshold {
		fmt.Printf("%s was caught!\n", *pokemonName)
		cfg.PokeDex[*pokemonName] = pokeData
	} else {
		fmt.Printf("%s escaped!\n", *pokemonName)
	}

	return nil
}

func CommandInspect(cfg *Config, pokemonName *string) error {
	pokeData, exists := cfg.PokeDex[*pokemonName]
	if !exists {
		fmt.Println("you have not caugh that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", *pokemonName)
	fmt.Printf("Height: %d\n", pokeData.Height)
	fmt.Printf("Weight: %d\n", pokeData.Weight)
	fmt.Printf("Stats:")
	for _, stat := range pokeData.Stats {
		fmt.Printf("\n -%s: %d", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("\nTypes:")
	for _, t := range pokeData.Types {
		fmt.Printf("\n -%s", t.Type.Name)
	}
	fmt.Println()

	return nil
}

func CommandPokedex(cfg *Config, _ *string) error {
	if len(cfg.PokeDex) == 0 {
		return fmt.Errorf("no pokemon caught yet")
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.PokeDex {
		fmt.Printf("- %s\n", name)
	}

	return nil
}
