package main

import (
	"log"

	"github.com/westside-jpg/WordMine/formatting"
	"github.com/westside-jpg/WordMine/services"
	"github.com/westside-jpg/WordMine/types"
)

func main() {
	list, err := services.TopWords("warandpeace.txt", 
	types.TopWordsOptions{
		Limit: 10,
		ExcludeStopWords: true,
	})

	if err != nil {
		log.Fatalf("Ошибка анализа слов по длине: %v", err)
	}

	formatting.TopWordsFormatting(list)

	list2, err := services.Stats("warandpeace.txt")
	formatting.StatsFormatting(list2)

	list3, err := services.FindInText("warandpeace.txt", types.FindInTextOptions{
		Words: []string{"мир"},
		WholeWordOnly: false,
		CaseSensitive: false,
	})
	formatting.FindInTextFormatting(list3)

	list4, err := services.LetterFrequency("warandpeace.txt")
	formatting.LetterFrequencyFormatting(list4)

	list5, err := services.TopNGrams("warandpeace.txt", types.TopNGramOptions{
		N: 2,
	})
	formatting.NGramsFormatting(list5)
}