package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	words := LoadWords("words.txt")

	for index, word := range words {
		fmt.Println(index, "--", word)
	}
	var response string
	test := GetWord(words)
	fmt.Printf("Try the following word: %s\n", test)
	fmt.Printf("What did you get? (b)lack, (y)ellow, (g)reen\n")
	// TODO check validity of input
	fmt.Scanln(&response)
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

func GetWord(words []string) string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(words))
	pick := words[randomIndex]
	return pick
}
