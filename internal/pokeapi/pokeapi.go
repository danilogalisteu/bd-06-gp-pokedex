package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PokeLocations struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokeEncounters struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func getUrl(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting url:\n%s", err)
		return nil
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("Error reading http body:\n%s", err)
		return nil
	}
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return nil
	}
	return body
}

var locations = PokeLocations{}
var encounters = PokeEncounters{}
var cache = NewCache(5 * time.Minute)

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
