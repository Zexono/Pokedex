package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main(){
	
	scanner := bufio.NewScanner(os.Stdin)
	con := &config{}
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		user_input := scanner.Text()
		command := cleanInput(user_input)

		if val,have := getCommands()[command[0]]; have {
			
			if val.name == "explore" && len(command) > 1{
				con.areaName = command[1]
				err := val.callback(con)
				if err != nil {
				fmt.Println(err)
				}
			}else if  val.name == "catch" && len(command) > 1{
				con.pokemonName = command[1]
				err := val.callback(con)
				if err != nil {
				fmt.Println(err)
				}
			}else {
				err := val.callback(con)
				if err != nil {
				fmt.Println(err)
				}
			}
		}else {
			commandNotFound(con)
		}

	}

}

func cleanInput(text string) []string{
	//var clean []string
	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	return words
}


