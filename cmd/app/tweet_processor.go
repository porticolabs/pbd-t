package main

import (
	"strings"
)

var (
	ignoreMarks = []string{".",",","!","?",":",";","\"","'","`"}
	ignoreWords = map[string]int{
		"el": 0, 
		"la": 0, 
		"los": 0, 
		"las": 0, 
		"un": 0, 
		"una": 0, 
		"unos": 0, 
		"unas": 0, 
		"a": 0, 
		"de": 0, 
		"para": 0, 
		"por": 0, 
		"con": 0, 
		"sin": 0, 
		"hola": 0, 
		"hey": 0, 
		"me": 0, 
		"en": 0, 
		"que": 0, 
		"tienen": 0, 
		"buen": 0,
		"bueno": 0, 
		"buena": 0, 
		"algo": 0, 
		"algún": 0}

	keywords = map[string]map[string]int{
		"Book": map[string]int{
			"libro": 3,
			"libros": 3,
			"leer": 1,
			"lectura": 1,
			"hojear": 1,
			"literatura": 2,
			},
		"Comic" : map[string]int{
			"comic": 3,
			"comics": 3,
			"leer": 1,
			"lectura": 1,
			"hojear": 1,
			},
		"Movie" : map[string]int{
			"ver": 1,
			"mirar": 1,
			"peli": 3,
			"película": 3,
			"pelicula": 3,
			"pelis": 3,
			"películas": 3,
			"peliculas": 3,
			"cine": 2,
			"largometraje": 2,
			},
		"Serie" : map[string]int{
			"ver": 1,
			"mirar": 1,
			"serie": 3,
			"series": 3,
			"capítulo": 2,
			"capitulo": 2,
			"episodio": 2,
			"temporada": 2,
			"capítulos": 2,
			"capitulos": 2,
			"episodios": 2,
			"temporadas": 2,
			},
	}
	keywordsTypeList = []string{"Book", "Movie"}
)

func prepareText(tweetText string) (string, map[string]int) {
	textIP := strings.ToLower(tweetText) // Text In Process
	for _, mark := range ignoreMarks {
		textIP = strings.ReplaceAll(textIP, mark, "")
	}
	textWords := strings.Fields(textIP) // Slice with the tweet text words

	// Each word will have a weight acording to the order in which they are
	//  in the sentence. This way, words that were written first, will weight
	//  more than the last words.
	wordsWeight := make(map[string]int)  // Map in which words weights will be stored
	wordIndex := len(textWords)+2 // Index weight 
	outputText := "" // Text with the final keywords to be analized

	for _, word := range textWords  {
		if _, ok := ignoreWords[word]; !ok {
			outputText += word + " "
			wordsWeight[word] = wordIndex
		}
		wordIndex--
    }
	return outputText, wordsWeight
}

func evaluateWords(keywordText string, keywordsWeight map[string]int, requestTypeKeywords map[string]int) int {
	matchCounter := 0
	for keyword, typeWeight  := range requestTypeKeywords {
		if strings.Contains(keywordText, keyword) {
			matchCounter += (typeWeight * keywordsWeight[keyword])
		}
	}
	return matchCounter
}

func analizeTweetText(tweetText string) (requestType []string) {
	// Analizing tweet
	tweetTextKeywords, wordsWeight := prepareText(tweetText)
	requestTypesFound := make(map[string]int)

	for requestType, requestTypeKeywords  := range keywords {
		typeEvaluation := evaluateWords(tweetTextKeywords, wordsWeight, requestTypeKeywords)
		if typeEvaluation > 0 {
			requestTypesFound[requestType] = typeEvaluation
		}
	}

	max := 0
	finalTypes := []string{}
	for requestType, requestTypeWeight := range requestTypesFound {
		if requestTypeWeight > max {
			finalTypes = []string{requestType}
			max = requestTypeWeight
		} else if requestTypeWeight == max {
			finalTypes = append(finalTypes, requestType)
		}
	}
	return finalTypes
}
