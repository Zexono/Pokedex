package main

import (
	"fmt"
	"os"
)

func commandNotFound() error{
	fmt.Println("Unknown command")
	return nil
}

func commandExit() error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	for _, val := range getCommands() {
		fmt.Printf("%s: %s",val.name,val.description)
		fmt.Println()
	}
	return nil
}


type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	return 	map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
		}
}