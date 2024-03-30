package main

import (
	"fmt"
	"net/http"
	"io"
)

type PokeLocation struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getUrl(url string) string {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting url:\n%s", err)
		return ""
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("Error reading http body:\n%s", err)
		return ""
	}
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return ""
	}
	return string(body)
}

func getLocations() {
	urlLocation := "https://pokeapi.co/api/v2/location/"
	locations := getUrl(urlLocation)
	fmt.Println(locations)
}
