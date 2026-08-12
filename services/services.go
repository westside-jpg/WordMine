package services

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// TODO: Сделать детерменированным
func GetTopCountByLength(name string, length int, offset int, caseSensetive bool) (map[string]int, error) {
	file, err := os.Open(name)
	if err != nil {
		log.Printf("Ошибка открытия файла: %v", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	rawList := make(map[string]int)
	for scanner.Scan() {
		word := scanner.Text()
		if !caseSensetive {
			word = strings.ToLower(word)
		}
		if len([]rune(word)) == length {
			rawList[word]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	list := make(map[string]int)
	for range offset {
		var maxCount = 0
		var maxCountWord string
		for word, count := range rawList {
			if count >= maxCount {
				maxCount = count
				maxCountWord = word
			}
		}
		list[maxCountWord] = maxCount
		delete(rawList, maxCountWord)
	}

	return list, nil
	
}