package main

import (
    "testing"
	"fmt"
	"strings"
)

func TestPrepareText1(t *testing.T){
	tweetText := "Hola @portico, me podés tirar algo bueno para ver? #quieroscifi"
	expectedResult := "@portico podés tirar ver #quieroscifi "

	tweetTextProcessed, tweetTextKeywords := prepareText(tweetText)

	fmt.Println(tweetTextProcessed)
	fmt.Println(tweetTextKeywords)

	if tweetTextProcessed != expectedResult {
		t.Fatalf("Expected text:\n\"%s\"\nbut got:\n\"%s\"", expectedResult, tweetTextProcessed)
	}

	if len(tweetTextKeywords) != 5 {
		t.Fatalf("Expected keywords:\n\"%d\"\nbut got:\n\"%d\"", 5, len(tweetTextKeywords))
	}
	
	if tweetTextKeywords["ver"] != 4 {
		t.Fatalf("Expected \"ver\" keyword value:\n\"%d\"\nbut got:\n\"%d\"", 4, tweetTextKeywords["ver"])
	}
}

func TestPrepareText2(t *testing.T){
	tweetText := "Hola @porticocba que película basada en un libro me recomiendan? #quieroscifi"
	expectedResult := "@porticocba película basada libro recomiendan #quieroscifi "

	tweetTextProcessed, tweetTextKeywords := prepareText(tweetText)

	if tweetTextProcessed != expectedResult {
		t.Fatalf("Expected text:\n\"%s\"\nbut got:\n\"%s\"", expectedResult, tweetTextProcessed)
	}

	if len(tweetTextKeywords) != 6 {
		t.Fatalf("Expected keywords:\n\"%d\"\nbut got:\n\"%d\"", 5, len(tweetTextKeywords))
	}
	
	if tweetTextKeywords["película"] != 10 {
		t.Fatalf("Expected \"ver\" keyword value:\n\"%d\"\nbut got:\n\"%d\"", 10, tweetTextKeywords["película"])
	}
}

func TestEvaluateWords1(t *testing.T){
	tweetText := "@porticocba #quieroscifi tirame copado leer"
	tweetWordsWeigths := map[string]int{
		"@porticocba": 7,
		"#quieroscifi": 6,
		"tirame": 5,
		"copado": 4,
		"leer": 3,
		}

	value := evaluateWords(tweetText, keywords["Book"], tweetWordsWeigths)
	if value != 3 {
		t.Fatalf("Expected type value:\n\"%d\"\nbut got:\n\"%d\"", 3, value)
	}
}

func TestEvaluateWords2(t *testing.T){
	tweetText := "Hola @porticocba que película basada en un libro me recomiendan para ver?"
	tweetWordsWeigths := map[string]int{
		"@porticocba": 12,
		"película": 10,
		"basada": 9,
		"libro": 6,
		"recomiendan": 4,
		"ver": 2,
		}

	value := evaluateWords(tweetText, keywords["Book"], tweetWordsWeigths)
	if value != 18 {
		t.Fatalf("Expected type value:\n\"%d\"\nbut got:\n\"%d\"", 18, value)
	}
}

func TestAnalizeTweetText1(t *testing.T){
	tweetText := "Hola @porticocba #quieroscifi que película basada en un libro me recomiendan para ver?"
	tweetRequestTypes := strings.Join(analizeTweetText(tweetText), ",")
	expectedTweetRequestTypes := "Movie"
	if tweetRequestTypes != expectedTweetRequestTypes {
		t.Fatalf("Expected types:\n\"%s\"\nbut got:\n\"%s\"", expectedTweetRequestTypes, tweetRequestTypes)
	}
}

func TestAnalizeTweetText2(t *testing.T){
	tweetText := "Hola @porticocba tengo ganas de leer un buen libro que haya tenido película? #quieroscifi"
	tweetRequestTypes := strings.Join(analizeTweetText(tweetText), ",")
	expectedTweetRequestTypes := "Book"
	if tweetRequestTypes != expectedTweetRequestTypes {
		t.Fatalf("Expected types:\n\"%s\"\nbut got:\n\"%s\"", expectedTweetRequestTypes, tweetRequestTypes)
	}
}