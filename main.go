package main

import (
	"log"

	"github.com/westside-jpg/WordMine/services"
	"github.com/westside-jpg/WordMine/utils"
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

	utils.GetTopFormatting(list)
}