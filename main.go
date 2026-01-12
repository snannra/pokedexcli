package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cfg := &config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}
	var commands map[string]cliCommand
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Lists locations in pokemon map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists previous locations in pokemon map",
			callback:    commandMapB,
		},
	}

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(cfg *config) error {
			fmt.Println("Usage:\n")
			for _, cmd := range commands {
				fmt.Printf("%s: %s\n", cmd.name, cmd.description)
			}
			return nil
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		command := cleanInput(input)[0]
		if cmd, exists := commands[command]; exists {
			err := cmd.callback(cfg)
			if err != nil {
				fmt.Printf("Error executing command %q: %v\n", command, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
