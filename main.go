package main

import (
	"log"

	"github.com/westside-jpg/WordMine/formatting"
	"github.com/westside-jpg/WordMine/services"
	"github.com/westside-jpg/WordMine/types"
)

func main() {
	list, err := services.Top("warandpeace.txt", 
	types.TopOptions{
		Limit: 10,
		ExcludeStopWords: true,
	})

	if err != nil {
		log.Fatalf("Ошибка анализа слов по длине: %v", err)
	}

	formatting.TopFormatting(list)

	list2, err := services.Stats("warandpeace.txt")
	formatting.StatsFormatting(list2)

	list3, err := services.FindInText("warandpeace.txt", types.FindInTextOptions{
		Words: []string{"король", "князь"},
		WholeWordOnly: false,
		CaseSensitive: true,
	})
	formatting.FindInTextFormatting(list3)

	list4, err := services.LetterFrequency("warandpeace.txt")
	formatting.LetterFrequencyFormatting(list4)

}