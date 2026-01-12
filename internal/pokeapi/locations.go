package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/snannra/pokedexcli/internal/pokeapi"
)

type locationAreaListResp struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func commandMap(cfg *pokeapi.config) error {
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

func commandMapB(cfg *config) error {
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
