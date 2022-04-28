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
	words_file := LoadWords("words.txt")

	var response string
	// Begin loop
	for response != "ggggg" {
		response = ""
		test := GetWord(words)
		fmt.Printf("Try the following word: %s\n", test)
		fmt.Printf("What did you get? (b)lack, (y)ellow, (g)reen\n")
		fmt.Scanln(&response)
		fmt.Printf("\nReceived: %s\n", response)

		valid := checkInputValidity(response)

		if len(response) > 0 && valid {
			words = FilterWords(words, test, response)
		}

		if response == "c" {
			words_file = remove(words_file, test)
			WriteWords("words.txt", words_file)
			fmt.Printf("Removed: %s\n", test)
		}

		if !valid {
			fmt.Println("ERROR: Invalid input!")
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

func WriteWords(path string, words []string) {
	log.Printf("WriteWords %s\n", path)

	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for _, word := range words {

		_, err := f.WriteString(word + "\n")

		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetWord(words []string) string {

	pick := "aa"
	randomIndex := -1

	for i := 0; i < len(words) && checkDoubles(pick); i++ {
		rand.Seed(time.Now().UnixNano())
		randomIndex = rand.Intn(len(words))
		pick = words[randomIndex]
	}

	return pick
}

func checkDoubles(word string) bool {

	occurrence := make(map[byte]bool, 0)

	for i := 0; i < len(word); i++ {
		c := word[i]
		_, ok := occurrence[c]

		if ok {
			return true
		} else {
			occurrence[c] = true
		}
	}

	return false
}

func FilterWords(toFilter []string, tested string, response string) []string {
	// TODO cosa succede parola corretta dodge e c'è c'ho la d centrale green su bodge? Mi cancella dodge?

	filtered := make([]string, 0)

	blacks := make([]byte, 0)
	others := make([]byte, 0)
	y_rexes := make([]string, 0)
	var g_rex string

	for i, r := range response {
		if r == 'b' {
			blacks = append(blacks, tested[i])
			g_rex = g_rex + "."
		} else if r == 'y' {
			others = append(others, tested[i])
			g_rex = g_rex + "."

			var y_rex string

			for j := 0; j < 5; j++ {
				if j == i {
					y_rex += string(tested[i])
				} else {
					y_rex += "."
				}
			}

			y_rexes = append(y_rexes, y_rex)

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

		// Remove all words that have yellow letters in the current position
		toSkip = false

		for _, rex := range y_rexes {
			if res, _ := regexp.MatchString(rex, word); res {
				toSkip = true
			}
		}

		if toSkip {
			continue
		}

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

func checkInputValidity(input string) bool {

	if len(input) != 5 {
		return false
	}

	if res, _ := regexp.MatchString("(g|b|y){5}", input); res {
		return true
	}

	return false
}

func remove(s []string, to_remove string) []string {
	i := -1
	for i = range s {
		if s[i] == to_remove {
			break
		}
	}

	if i < 0 {
		return s
	}

	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
