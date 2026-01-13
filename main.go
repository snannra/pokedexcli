package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/snannra/pokedexcli/internal/pokecache"
	"github.com/snannra/pokedexcli/internal/repl"
)

func main() {
	cfg := &repl.Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		Cache:    pokecache.NewCache(5 * time.Second),
		PokeDex:  make(map[string]repl.Pokemon),
	}
	var commands map[string]repl.CliCommand
	commands = map[string]repl.CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    repl.CommandExit,
		},
		"map": {
			Name:        "map",
			Description: "Lists locations in pokemon map",
			Callback:    repl.CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Lists previous locations in pokemon map",
			Callback:    repl.CommandMapB,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore a location area to see which Pokémon can be found there",
			Callback:    repl.CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catch a pokemon by choosing it's name",
			Callback:    repl.CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a caught pokemon by choosing it's name",
			Callback:    repl.CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Lists all caught pokemon in your pokedex",
			Callback:    repl.CommandPokedex,
		},
	}

	commands["help"] = repl.CliCommand{
		Name:        "help",
		Description: "Displays a help message",
		Callback: func(cfg *repl.Config, _ *string) error {
			fmt.Println("Usage:")
			for _, cmd := range commands {
				fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
			}
			return nil
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := repl.CleanInput(input)
		command := cleanedInput[0]
		var argTwo string
		if len(cleanedInput) > 1 {
			argTwo = cleanedInput[1]
		}
		if cmd, exists := commands[command]; exists {
			err := cmd.Callback(cfg, &argTwo)
			if err != nil {
				fmt.Printf("Error executing command %q: %v\n", command, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
