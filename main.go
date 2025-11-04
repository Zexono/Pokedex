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
		if val,have := user_command[command[0]]; have {
			//fmt.Println(val.description)
			val.callback()
		}else {
			commandNotFound()
		}
		//fmt.Println("Your command was: "+command[0])
		
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
	for _, val := range user_command_for_help {
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

var user_command = map[string]cliCommand{
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

var user_command_for_help = map[string]cliCommand{
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
    },
	"help": {
        name:        "help",
        description: "Displays a help message",
	},
}