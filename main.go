package main

import (
	"fmt"
	"strings"
)
func main(){
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string{
	//var clean []string
	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	return words
}
