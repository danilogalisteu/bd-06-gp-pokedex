package main

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
var cache = NewCache(5 * time.Minute)

func getLocationsNext() PokeLocations {
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

func getLocationsPrev() PokeLocations {
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
