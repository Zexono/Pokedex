package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Zexono/pokedexcli/internal"
)

var cache = internal.NewCache(20 * time.Second)
//internal.Cache{}

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

func commandMap() error {
    var url string
    if location.Next != "" {
        url = location.Next
    } else {
        url = "https://pokeapi.co/api/v2/location-area/"
    }

    var body []byte
    if data, ok := cache.Get(url); ok {
		fmt.Println("[cache] hit:", url)
        body = data
    } else {
        res, err := http.Get(url)
		fmt.Println("[cache] miss:", url)
        if err != nil {
            return err
        }
        defer res.Body.Close()

        b, err := io.ReadAll(res.Body)
        if err != nil {
            return err
        }
        if res.StatusCode > 299 {
            return fmt.Errorf("status %d: %s", res.StatusCode, string(b))
        }
        cache.Add(url, b)
        body = b
    }

    if err := json.Unmarshal(body, &location); err != nil {
        return err
    }
    for _, loc := range location.Results {
        fmt.Println(loc.Name)
    }
    return nil
}

func commandMapBack() error{
	/*if location.Previous == ""{
		fmt.Println("no previous map from here")
		return nil
	}*/

	//if location.Previous == "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20" || 
	if location.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	url := location.Previous

	//var res *http.Response
	//var err error
	var body []byte
    if data, ok := cache.Get(url); ok {
		fmt.Println("[cache] hit:", url)
        body = data
    }else{
		res, err :=  http.Get(url)
		fmt.Println("[cache] miss:", url)
		if err != nil{
			return  err
		}

		b, err := io.ReadAll(res.Body)
		defer res.Body.Close()

		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil{
			return  err
		}
		cache.Add(url, b)
        body = b
	}
	err := json.Unmarshal(body, &location)
	if err != nil{
		return  err
	}
	
	for _, loc := range location.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
	confiq		*Config
}

type Config struct {
	Next 		string `json:"next"`
	Previous 	string `json:"previous"`
	Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}

var location Config

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
			"map": {
				name:        "map",
				description: "Displays a next 20 location",
				callback:    commandMap,
				confiq: 	 &Config{},
			},
			"mapb": {
				name:        "mapb",
				description: "Displays a previous 20 location",
				callback:    commandMapBack,
				confiq: 	 &Config{},
			},
		}
}

