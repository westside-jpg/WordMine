package services

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/westside-jpg/WordMine/stopwords"
	"github.com/westside-jpg/WordMine/types"

	"github.com/kljensen/snowball"
)

// GetTop читает текстовый файл по указанному пути и возвращает топ самых
// часто встречающихся слов, отсортированных по убыванию частоты. При равной
// частоте слова упорядочиваются по алфавиту, чтобы результат был одинаковым
// при каждом вызове.
//
// Слова очищаются от знаков препинания по краям перед подсчётом (например,
// "привет," и "привет" считаются одним словом), но внутренние символы,
// такие как дефис в "кто-то", сохраняются.
//
// Параметры:
//   - name: путь к текстовому файлу для анализа
//   - options: настройки подсчёта, см. types.TopOptions
//     Length: если 0, учитываются слова любой длины; иначе только слова
//     ровно этой длины в рунах
//     Limit: сколько слов вернуть в топе; если 0, используется значение
//     по умолчанию 10
//     CaseSensitive: если false, регистр букв игнорируется при подсчёте
//     ExcludeStopWords: если true, из результата исключаются предлоги,
//     союзы, частицы и местоимения. Список берётся из StopWords, если
//     он задан, иначе используется встроенный stopwords.DefaultRussianStopWords.
//     Проверка на стоп-слово регистронезависима всегда, вне зависимости
//     от CaseSensitive
//     StopWords: свой список стоп-слов, учитывается только если
//     ExcludeStopWords равен true
//
// Возвращает ошибку, если файл не удалось открыть или прочитать.
// Если найдено меньше слов, чем запрошено в Limit, возвращает столько,
// сколько нашлось, без ошибки.
func GetTop(name string, options types.TopOptions) ([]types.WordCount, error) {
	if options.Limit == 0 {
		options.Limit = 10
	}

	stopWords := options.StopWords
	if options.ExcludeStopWords && stopWords == nil {
		stopWords = stopwords.DefaultRussianStopWords
	}

	stemmedStopWords := make(map[string]struct{}, len(stopWords))
	for _, sw := range stopWords {
		stemmed, err := snowball.Stem(strings.ToLower(sw), "russian", true)

		if err != nil {
			return []types.WordCount{}, err
		}
		
		stemmedStopWords[stemmed] = struct{}{}
	}

	file, err := os.Open(name)
	if err != nil {
		log.Printf("Ошибка открытия файла: %v", err)
		return []types.WordCount{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	rawList := make(map[string]int)
	for scanner.Scan() {
		word := scanner.Text()

		word = strings.TrimFunc(word, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if word == "" {
			continue
		}

		if !options.CaseSensitive {
			word = strings.ToLower(word)
		}

		if options.ExcludeStopWords {
			stemmedWord, err := snowball.Stem(word, "russian", true)
			if err != nil {
				return []types.WordCount{}, err
			}
			if _, isStopWord := stemmedStopWords[strings.ToLower(stemmedWord)]; isStopWord {
				continue
			}
		}

		if options.Length == 0 || len([]rune(word)) == options.Length {
			rawList[word]++
		}
	}

	if err := scanner.Err(); err != nil {
		return []types.WordCount{}, err
	}
	
	results := make([]types.WordCount, 0, len(rawList))
	for word, count := range rawList {
		results = append(results, types.WordCount{Word: word, Count: count})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Word < results[j].Word
	})

	if options.Limit > len(results) {
		options.Limit = len(results)
	}

	return results[:options.Limit], nil
	
}