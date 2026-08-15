package main

import (
	"log"

	"github.com/westside-jpg/WordMine/formatting"
	"github.com/westside-jpg/WordMine/services"
	"github.com/westside-jpg/WordMine/types"
)

func main() {
	list, err := services.GetTop("warandpeace.txt", 
	types.TopOptions{
		Limit: 10,
		ExcludeStopWords: true,
	})

	if err != nil {
		log.Fatalf("Ошибка анализа слов по длине: %v", err)
	}

	formatting.GetTopFormatting(list)

	list2, err := services.GetStats("warandpeace.txt")
	formatting.GetStatsFormatting(list2)

	list3, err := services.FindInText("warandpeace.txt", types.FindInTextOptions{
		Words: []string{"евро"},
		WholeWordOnly: false,
		CaseSensitive: false,
	})
	formatting.FindInTextFormatting(list3)

}