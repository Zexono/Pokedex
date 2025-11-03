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
		fmt.Println("Your command was: "+command[0])
		
	}
}

func cleanInput(text string) []string{
	//var clean []string
	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	return words
}
