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

	var response string
	// Begin loop
	for response != "ggggg" {
		response = ""
		test := GetWord(words)
		fmt.Printf("Try the following word: %s\n", test)
		fmt.Printf("What did you get? (b)lack, (y)ellow, (g)reen\n")
		// TODO check validity of input
		fmt.Scanln(&response)
		fmt.Printf("\nReceived: %s", response)

		if len(response) > 0 {
			words = FilterWords(words, test, response)
		}
	}
}

func LoadWords(path string) []string {
	// Load word seed keeps ones with lenght 5, lowers all of them and keeps the ones matching [a-z]{5}
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

func FilterWords(toFilter []string, tested string, response string) []string {
	filtered := make([]string, 0)

	blacks := make([]byte, 0)
	others := make([]byte, 0)
	var g_rex string

	for i, r := range response {
		if r == 'b' {
			blacks = append(blacks, tested[i])
			g_rex = g_rex + "."
		} else if r == 'y' {
			others = append(others, tested[i])
			g_rex = g_rex + "."
		} else {
			others = append(others, tested[i])
			g_rex = g_rex + string(tested[i])
		}
	}

	for _, word := range toFilter {
		// If contains a black letter; continue
		toSkip := false
		for _, r := range blacks {

			if contains(others, r) {
				// In this scenario this letter appears at least twice in the word
				// At the time being I skip it when black, however there is information here to exploit in the future.
				// Probably, I just have to remove all the words that have this letter in this position
				continue
			}

			if strings.Contains(word, string(r)) {
				toSkip = true
			}
		}

		if toSkip {
			continue
		}
		// If not contains all green and yellow letters; continue
		toSkip = false
		for _, r := range others {
			if !strings.Contains(word, string(r)) {
				toSkip = true
			}
		}

		if toSkip {
			continue
		}
		// If not matches regex about green letters in given position; continue
		// eg. bagel and ge is corrent, regex -> ..ge.

		if res, _ := regexp.MatchString(g_rex, word); !res {
			continue
		}

		// TODO remove all words that have yellow letters in the current position
		filtered = append(filtered, word)
	}

	log.Printf("Filtered words: %d -> %d", len(toFilter), len(filtered))
	return filtered
}

func contains(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
