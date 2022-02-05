package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	words := LoadWords("words.txt")

	for index, word := range words {
		fmt.Println(index, "--", word)
	}
}

// LoadWords(file path) -> list of valid words
func LoadWords(path string) []string {
	// Load word seed and remove all words !=5, lower all them and keep the ones matching [a-z]{5}
	log.Printf("LoadWords %s\n", path)

	words := make([]string, 0)

	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		eval := strings.ToLower(scanner.Text())

		if len(eval) == 5 {
			if res, _ := regexp.MatchString("[a-z]{5}", eval); res {
				words = append(words, eval)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return words
}
