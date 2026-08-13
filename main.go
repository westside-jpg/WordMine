package main

import (
	"fmt"
	"log"

	"github.com/westside-jpg/WordMine/services"
	"github.com/westside-jpg/WordMine/utils"
	"github.com/westside-jpg/WordMine/types"
	
)

func main() {
	list, err := services.GetTop("warandpeace.txt", 
	types.TopOptions{
		ExcludeStopWords: true,
	})

	if err != nil {
		log.Fatalf("Ошибка анализа слов по длине: %v", err)
	}

	for _, s := range list {
		fmt.Printf("Слово \"%s\" встретилось ровно %d %s\n", s.Word, s.Count, utils.DeclinationWord(s.Count, "раз", "раза", "раз"))
	}
}