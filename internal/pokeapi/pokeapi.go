package pokeapi

import (
	"encoding/json"
	"fmt"
	"time"
)

var cache = NewCache(5 * time.Minute)
var locations = PokeLocations{}
var encounters = PokeEncounters{}
var info = PokeInfo{}

func GetLocationsNext() PokeLocations {
	urlLocation := "https://pokeapi.co/api/v2/location-area/"
	if locations.Next != "" {
		urlLocation = locations.Next
	}

	if _, ok := cache.Get(urlLocation); !ok {
		content := getUrl(urlLocation)
		cache.Add(urlLocation, content)
	}
	content, _ := cache.Get(urlLocation)

	err := json.Unmarshal(content, &locations)
	if err != nil {
		fmt.Println(err)
	}

	return locations
}

func GetLocationsPrev() PokeLocations {
	if locations.Previous == nil {
		fmt.Println("No previous page of results")
		return locations
	}
	urlLocation := *locations.Previous

	if _, ok := cache.Get(urlLocation); !ok {
		content := getUrl(urlLocation)
		cache.Add(urlLocation, content)
	}
	content, _ := cache.Get(urlLocation)

	err := json.Unmarshal(content, &locations)
	if err != nil {
		fmt.Println(err)
	}

	return locations
}

func ExploreLocation(id string) (PokeEncounters, bool) {
	urlLocation := ""
	for _, res := range locations.Results {
		if id == res.Name {
			urlLocation = res.URL
			break
		}
	}

	if urlLocation == "" {
		return encounters, false
	}

	if _, ok := cache.Get(urlLocation); !ok {
		content := getUrl(urlLocation)
		cache.Add(urlLocation, content)
	}
	content, _ := cache.Get(urlLocation)

	err := json.Unmarshal(content, &encounters)
	if err != nil {
		fmt.Println(err)
	}

	return encounters, true
}

func CatchPokemon(id string) (PokeInfo, bool) {
	urlLocation := ""
	for _, encounter := range encounters.PokemonEncounters {
		if id == encounter.Pokemon.Name {
			urlLocation = encounter.Pokemon.URL
			break
		}
	}

	if urlLocation == "" {
		return info, false
	}

	if _, ok := cache.Get(urlLocation); !ok {
		content := getUrl(urlLocation)
		cache.Add(urlLocation, content)
	}
	content, _ := cache.Get(urlLocation)

	err := json.Unmarshal(content, &info)
	if err != nil {
		fmt.Println(err)
	}

	return info, true
}
