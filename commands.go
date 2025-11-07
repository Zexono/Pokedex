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
const(
	baseURL = "https://pokeapi.co/api/v2/location-area/"
)
var cache = internal.NewCache(20 * time.Second)
//internal.Cache{}

func commandNotFound(_ *config) error{
	fmt.Println("Unknown command")
	return nil
}

func commandExit(_ *config) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	for _, val := range getCommands() {
		fmt.Printf("%s: %s",val.name,val.description)
		fmt.Println()
	}
	return nil
}

var location listLocation

func commandMap(_ *config) error {
    var url string
    if location.Next != "" {
        url = location.Next
    } else {
        url = baseURL//"https://pokeapi.co/api/v2/location-area/"
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

func commandMapBack(_ *config) error{

	if location.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	url := location.Previous

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

func commandExplore(con *config) error{
	var url string
	var area LocationArea
	if con.areaName == "" {
		return fmt.Errorf("plz input Area name you want to explore")
	}
	url = baseURL+con.areaName

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

    if err := json.Unmarshal(body, &area); err != nil {
        return err
    }
	
	fmt.Printf("Exploring %s... \n",con.areaName)
	fmt.Println("Found Pokemon: ")
    for _, pokemon := range area.PokemonEncounters {
        fmt.Printf(" - %s \n",pokemon.Pokemon.Name)
    }


	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
	
}

type config struct {
	areaName 	string
	//pokemonName	string //future if want pokemon stat
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
			"map": {
				name:        "map",
				description: "Displays a next 20 location",
				callback:    commandMap,	
			},
			"mapb": {
				name:        "mapb",
				description: "Displays a previous 20 location",
				callback:    commandMapBack,
			},
			"explore": {
				name:        "explore",
				description: "explore <area_name>\n	 Displays a Pokemon in area",
				callback:    commandExplore,
			},
		}
}

