package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/snannra/pokedexcli/internal/repl"
)

func main() {
	cfg := &repl.Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
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
	}

	commands["help"] = repl.CliCommand{
		Name:        "help",
		Description: "Displays a help message",
		Callback: func(cfg *repl.Config) error {
			fmt.Println("Usage:\n")
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
		command := repl.CleanInput(input)[0]
		if cmd, exists := commands[command]; exists {
			err := cmd.Callback(cfg)
			if err != nil {
				fmt.Printf("Error executing command %q: %v\n", command, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
