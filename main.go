package main

import (
	"fmt"
	"log"

	"github.com/westside-jpg/WordMine/services"
	"github.com/westside-jpg/WordMine/utils"
)

func main() {
	list := make(map[string]int)
	list, err := services.GetTopCountByLength("warandpeace.txt", 3, 10, false)

	if err != nil {
		log.Fatalf("Ошибка анализа слов по длине: %v", err)
	}

	for word, count := range list {
		fmt.Printf("Слово \"%s\" встретилось ровно %d %s\n", word, count, utils.DeclinationWord(count, "раз", "раза", "раз"))
	}
}