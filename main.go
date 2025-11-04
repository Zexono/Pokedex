package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main(){
	
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		user_input := scanner.Text()
		command := cleanInput(user_input)

		if val,have := getCommands()[command[0]]; have {
			err := val.callback()
			if err != nil {
				fmt.Println(err)
			}
		}else {
			commandNotFound()
		}

	}

}

func cleanInput(text string) []string{
	//var clean []string
	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	return words
}

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
