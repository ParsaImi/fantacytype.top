package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

type SentenceResponse struct {
	Sentence string `json:"sentence"`
	Locale   string `json:"locale"`
}

func GenerateCompetitionText(room, language string) (string, []string ) {
	
	
	// Common Persian words for room3
	persianWords := []string{
		"سلام", "جهان", "کتاب", "خانه", "درخت", "آزادی", "عشق", "دوست", "خورشید", "ماه",
		"ستاره", "آسمان", "زمین", "رود", "کوه", "گل", "پرنده", "باغ", "شهر", "روستا",
	}

	// English words for other rooms
	englishWords := []string{
		"introductory", "preliminary", "ceremonial", "demonstrating", "graves", "speed", "hoses", "steep", "reunions", "I", "You", "Yeah",
		"script", "devoted", "prepositions", "indie", "fascinating", "courage", "Star", "Five", "outside",
	}
	
	baseurl := "http://127.0.0.1:8000/corpus/sentence"

	u , _ := url.Parse(baseurl)

	q := u.Query()
	if language == "fa"{
		q.Set("locale", "fa_IR")
	}else{
		q.Set("locale", "en_US")
	}
	u.RawQuery = q.Encode()
	resp, _ := http.Get(u.String())
	defer resp.Body.Close()

	body, _:= io.ReadAll(resp.Body)

	var response SentenceResponse
	err := json.Unmarshal(body , &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}
	splitedWords := strings.Split(response.Sentence, " ")


	var words []string
	var result []string

	if room == "room3" {
		words = persianWords
	} else {
		words = englishWords
	}

	// Check if the word list is empty
	if len(words) == 0 {
		fmt.Println("ohhhahahahahha")
	}

	// Generate 10 random words
	for i := 0; i < 10; i++ {
		word := words[rand.Intn(len(words))]
		result = append(result, word)
	}


	return response.Sentence, splitedWords
}
