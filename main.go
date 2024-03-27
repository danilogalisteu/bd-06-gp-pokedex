package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Pokedex > ")
        scanner.Scan()
        input := scanner.Text()
		eparse := parseCommand(input)
		if eparse != nil {
			return
		}
	}
}
