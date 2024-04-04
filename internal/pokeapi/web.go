package pokeapi

import (
	"fmt"
	"io"
	"net/http"
)

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
