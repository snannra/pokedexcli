package repl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Next     string
	Previous string
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
}

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func CommandExit(cfg *Config) error {
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

func CommandMap(cfg *Config) error {
	const baseLocationAreaURL = "https://pokeapi.co/api/v2/location-area/"

	url := cfg.Next
	if url == "" {
		url = baseLocationAreaURL
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

func CommandMapB(cfg *Config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	url := cfg.Previous

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
