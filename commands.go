package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Zexono/pokedexcli/internal"
)
const(
	baseURL = "https://pokeapi.co/api/v2/location-area/"
	pokeURL = "https://pokeapi.co/api/v2/pokemon/"
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

func commandCatch(con *config) error{
	rand.Float32()
	var url string
	var pokemon pokemonData
	if con.pokemonName == "" {
		return fmt.Errorf("plz input Pokemon name you want to catch")
	}
	url = pokeURL+con.pokemonName

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

    if err := json.Unmarshal(body, &pokemon); err != nil {
        return err
    }
	
	fmt.Printf("Throwing a Pokeball at %s... \n",pokemon.Name)
	//edit catch difficulity later
	//difficulity := pokemon.BaseExperience
	//formular := 1 + (pokemon.BaseExperience - 1) * (99) / (325 - 1)
	difficulity := int(math.Round(float64(1 + (pokemon.BaseExperience - 1) * (99) / (325 - 1))))

	if rand.Intn(100) >= difficulity {
		key := pokemon.Name //what if want tp catch same pokemon? //edit later
		myPokedex[key] = pokemon
		fmt.Printf("%s was caught! \n",pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		return nil

	}

	fmt.Printf("%s escaped! \n",pokemon.Name)
	
	return nil
}

func commandInspect(con *config) error{
	if val,have := myPokedex[con.pokemonName] ; have{
		fmt.Printf("Name: %s \n",val.Name)
		fmt.Printf("Height: %v \n",val.Height)
		fmt.Printf("Weight: %v \n",val.Weight)
		//fmt.Printf("stats : %v \n",val.Stats)
		fmt.Println("Stats: ")
		for _, v := range val.Stats {
			fmt.Printf(" -%s: %v\n",v.Stat.Name,v.BaseStat)
		}
		fmt.Println("Type: ")
		for _, v := range val.Types {
			fmt.Printf(" -%s\n",v.Type.Name)
		}
		return nil
		
	}
	fmt.Println("you have not caught that pokemon")
	return nil
}

func commandPokedex(con *config) error{
	if len(myPokedex) == 0 {
		fmt.Println("Don't have any Pokemon yet")
		return nil
	}
	for _, v := range myPokedex {
		fmt.Printf(" -%s \n",v.Name)
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
	pokemonName	string //future if want pokemon stat
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
			"catch": {
				name:        "catch",
				description: "catch <pokemon_name>\n	 to catch a Pokemon",
				callback:    commandCatch,
			},
			"inspect": {
				name:        "inspect",
				description: "inspect <pokemon_name>\n	 allow players to see details about a Pokemon if they have seen it before (or in our case, caught it)",
				callback:    commandInspect,
			},
			"pokedex": {
				name:        "pokedex",
				description: "show a list of all the names of the Pokemon the user has caught",
				callback:    commandPokedex,
			},
		}
}

