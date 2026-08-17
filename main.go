package main

import (
	"github.com/westside-jpg/WordMine/formatting"
	"github.com/westside-jpg/WordMine/services"
	"github.com/westside-jpg/WordMine/types"
	"github.com/westside-jpg/WordMine/export"
)

func main() {
	list1, _ := services.TopWords("warandpeace.txt", 
	types.TopWordsOptions{
		Limit: 10,
		ExcludeStopWords: true,
	})
	// formatting.TopWordsFormatting(list1)
	export.ExportJSON(list1, "jsons/top_words.json")
	export.ExportTXT("txts/top_words.txt", func() { formatting.TopWordsFormatting(list1) })

	list2, _ := services.Stats("warandpeace.txt")
	// formatting.StatsFormatting(list2)
	export.ExportJSON(list2, "jsons/stats.json")
	export.ExportTXT("txts/stats.txt", func() { formatting.StatsFormatting(list2) })
	

	list3, _ := services.FindInText("warandpeace.txt", types.FindInTextOptions{
		Words: []string{"мир"},
		WholeWordOnly: false,
		CaseSensitive: false,
	})
	// formatting.FindInTextFormatting(list3)
	export.ExportJSON(list3, "jsons/find_in_text.json")
	export.ExportTXT("txts/find_in_text.txt", func() { formatting.FindInTextFormatting(list3) })

	list4, _ := services.LetterFrequency("warandpeace.txt")
	// formatting.LetterFrequencyFormatting(list4)
	export.ExportJSON(list4, "jsons/letters_freq.json")
	export.ExportTXT("txts/letters_freq.txt", func() { formatting.LetterFrequencyFormatting(list4) })


	list5, _ := services.TopNGrams("warandpeace.txt", types.TopNGramOptions{
		N: 2,
	})
	// formatting.NGramsFormatting(list5)
	export.ExportJSON(list5, "jsons/n_grams.json")
	export.ExportTXT("txts/n_grams.txt", func() { formatting.NGramsFormatting(list5) })
}