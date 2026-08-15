package services

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/westside-jpg/WordMine/stopwords"
	"github.com/westside-jpg/WordMine/types"
	"github.com/westside-jpg/WordMine/utils"

	"github.com/kljensen/snowball"
)

// Top читает текстовый файл по указанному пути и возвращает топ самых
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
func Top(name string, options types.TopOptions) ([]types.WordCount, error) {
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

		word = utils.CleanWord(word)

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

// Stats читает текстовый файл целиком и возвращает набор статистических
// и лингвистических характеристик текста: количество букв, слов,
// предложений, знаков препинания, цифр и чисел, лексическое разнообразие,
// долю стоп-слов и приблизительный индекс читаемости.
//
// Файл читается целиком в память через os.ReadFile,
// а не потоково, поскольку разбиение на предложения требует видеть весь
// текст сразу.
//
// Ряд метрик являются приближёнными, а не точными:
//   - границы предложений определяются по знакам .!? и не распознают
//     сокращения (например, "т.д.", "А.С. Пушкин") или десятичные дроби;
//   - число слогов (AvgWordLengthBySyllables) считается как число гласных
//     букв в слове, без учёта реальных фонетических правил;
//   - ReadabilityScore — адаптация формулы Флеша для русского языка
//     (коэффициенты Флеша-Оборневой), эвристическая оценка, а не точный
//     лингвистический показатель.
//
// Параметры:
//   - name: путь к текстовому файлу для анализа
//
// Возвращает ошибку, если файл не удалось открыть или прочитать, если
// текст не содержит ни одного предложения, или если после очистки в
// тексте не осталось ни одного слова.
func Stats(name string) (types.Stats, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return types.Stats{}, err
	}
	text := string(data)

	// Знаки пунктуации
	var totalPunctuation int

	// Буквы
	var totalSymbols int
	var totalSymbolsWithoutSpaces int
	var onlyWordsLetters int

	var totalFigures int
	var totalNumbers int

	// Слова
	var totalWords int
	var totalStopWords int
	var stopWordsPercentage float64
	var uniqueWords int
	var totalSyllables int
	var avgWordLengthByLetters float64
	var avgWordLengthBySyllables float64
	var longestWord string
	var shortestWord string

	// Предложения
	var totalSentences int
	var avgSentenceLength float64
	var longestSentence string
	var shortestSentence string

	// Лингвистические характеристики
	var typeTokenRatio float64
	var readabilityScore float64

	for _, r := range text {
		totalSymbols++
		if !unicode.IsSpace(r) {
			totalSymbolsWithoutSpaces++
		}

		if unicode.IsNumber(r) {
			totalFigures++
		}
	}

	sentences := utils.SplitIntoSentences(text)

	if len(sentences) == 0 {
		return types.Stats{}, fmt.Errorf("Текст не содержит предложений для анализа")
	}

	sortedSentences := slices.Clone(sentences)
	slices.SortFunc(sortedSentences, func(a, b string) int {
		return len([]rune(a)) - len([]rune(b))
	})

	totalSentences = len(sentences)
	shortestSentence = sortedSentences[0]
	longestSentence = sortedSentences[len(sortedSentences)-1]

	stemmedStopWords := make(map[string]struct{}, len(stopwords.DefaultRussianStopWords))
	for _, sw := range stopwords.DefaultRussianStopWords {
		stemmed, err := snowball.Stem(strings.ToLower(sw), "russian", true)
		if err != nil {
			return types.Stats{}, err
		}
		stemmedStopWords[stemmed] = struct{}{}
	}

	seenWords := make(map[string]struct{})
	for _, sentence := range sentences {

		words := strings.Fields(sentence)

		for _, word := range words {

			word = utils.CleanWord(word)

			if word == "" {
				continue
			}

			totalWords++
			totalSyllables += utils.CountSyllables(word)
			onlyWordsLetters += len([]rune(word))

			seenWords[strings.ToLower(word)] = struct{}{}

			if longestWord == "" || len([]rune(word)) > len([]rune(longestWord)) {
				longestWord = word
			}
			if shortestWord == "" || len([]rune(word)) < len([]rune(shortestWord)) {
				shortestWord = word
			}

			_, err := strconv.ParseFloat(word, 64)
			if err == nil {
				totalNumbers++
			}

			stemmedWord, err := snowball.Stem(strings.ToLower(word), "russian", true)
			if err != nil {
				return types.Stats{}, err
			}
			if _, isStopWord := stemmedStopWords[stemmedWord]; isStopWord {
				totalStopWords++
			}

		}
	}

	if totalWords == 0 {
		return types.Stats{}, fmt.Errorf("Текст не содержит слов для анализа")
	}

	uniqueWords = len(seenWords)
	avgSentenceLength = float64(totalWords) / float64(len(sentences))
	avgWordLengthByLetters = float64(onlyWordsLetters) / float64(totalWords)
	avgWordLengthBySyllables = float64(totalSyllables) / float64(totalWords)
	totalPunctuation = totalSymbolsWithoutSpaces - onlyWordsLetters
	stopWordsPercentage = (float64(totalStopWords) / float64(totalWords)) * 100

	typeTokenRatio = (float64(uniqueWords) / float64(totalWords)) * 100
	readabilityScore = 206.835 - 1.52*avgSentenceLength - 65.14*avgWordLengthBySyllables

	return types.Stats{
		TotalPunctuation: totalPunctuation,

		TotalSymbols:              totalSymbols,
		TotalSymbolsWithoutSpaces: totalSymbolsWithoutSpaces,
		OnlyWordsLetters:          onlyWordsLetters,

		TotalFigures: totalFigures,
		TotalNumbers: totalNumbers,

		TotalWords:               totalWords,
		TotalSyllables:           totalSyllables,
		TotalStopWords:           totalStopWords,
		StopWordsPercentage:      stopWordsPercentage,
		UniqueWords:              uniqueWords,
		AvgWordLengthByLetters:   avgWordLengthByLetters,
		AvgWordLengthBySyllables: avgWordLengthBySyllables,
		LongestWord:              longestWord,
		ShortestWord:             shortestWord,

		TotalSentences:    totalSentences,
		AvgSentenceLength: avgSentenceLength,
		LongestSentence:   longestSentence,
		ShortestSentence:  shortestSentence,

		TypeTokenRatio:   typeTokenRatio,
		ReadabilityScore: readabilityScore,
	}, nil
}

// FindInText читает текстовый файл по указанному пути построчно и ищет
// вхождения одного или нескольких слов из options.Words, возвращая для
// каждого найденного вхождения номер строки, позицию символа в строке,
// найденное слово, контекст вокруг него и то поисковое слово, которое
// вызвало совпадение.
//
// Поисковые слова из options.Words перед сравнением очищаются от
// пунктуации по краям (см. utils.CleanWord), приводятся к нижнему
// регистру при options.CaseSensitive равном false и дедуплицируются.
// Исходный слайс options.Words при этом не изменяется, вся подготовка
// идёт в локальную копию.
//
// Слово из текста считается совпавшим одним из двух способов, в
// зависимости от options.WholeWordOnly:
//   - true: очищенное слово из текста должно точно совпадать с поисковым
//     словом (WholeWordOnly).
//   - false: очищенное слово из текста должно содержать поисковое слово
//     как подстроку (strings.Contains).
//
// Если одно слово из текста совпадает сразу с несколькими поисковыми
// словами (например, при options.WholeWordOnly равном false, слово
// "европейского" совпадает и с "евро", и с "европ"), для него создаётся
// отдельная запись результата на каждое совпавшее поисковое слово,
// поле MatchedWord в каждой из них указывает, какое именно.
//
// Context строится из трёх соседних токенов строки (предыдущий, само
// слово, следующий) в их исходном виде, с сохранением пунктуации, а не
// из очищенных слов, чтобы контекст выглядел ближе к оригинальному
// тексту. На границах строки (первый или последний токен) отсутствующий
// сосед просто опускается.
//
// Результаты сортируются по возрастанию LineIndex, а внутри одной строки
// по возрастанию CharIndex, то есть в порядке появления в тексте.
//
// Параметры:
//   - name: путь к текстовому файлу для анализа.
//   - options: настройки поиска, см. types.FindInTextOptions.
//
// Возвращает ошибку, если файл не удалось открыть или прочитать. Если
// совпадений не найдено, возвращает пустой Results без ошибки.
func FindInText(name string, options types.FindInTextOptions) (types.FindInTextResponse, error) {
	file, err := os.Open(name)
	if err != nil {
		return types.FindInTextResponse{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var results []types.FindInTextResults
	lineNumber := 0

	searchWords := make([]string, 0, len(options.Words))
	for _, w := range options.Words {
		cleaned := utils.CleanWord(w)
		if !options.CaseSensitive {
			cleaned = strings.ToLower(cleaned)
		}
		if cleaned != "" {
			searchWords = append(searchWords, cleaned)
		}
	}

	// Дедупликация повторений
	slices.Sort(searchWords)
	searchWords = slices.Compact(searchWords)

	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		tokens := utils.SplitLineWithPositions(line)

		for n, token := range tokens {
			word := utils.CleanWord(token.Word)
			if word == "" {
				continue
			}

			sameWord := word

			if !options.CaseSensitive {
				sameWord = strings.ToLower(word)
			}

			for _, sw := range searchWords {

				if (sameWord == sw && options.WholeWordOnly) || (strings.Contains(sameWord, sw) && !options.WholeWordOnly) {
					var context string

					switch {
					case n > 0 && n < len(tokens)-1:
						context = tokens[n-1].Word + " " + token.Word + " " + tokens[n+1].Word
					case n > 0:
						context = tokens[n-1].Word + " " + token.Word
					case n < len(tokens)-1:
						context = token.Word + " " + tokens[n+1].Word
					default:
						context = token.Word
					}

					results = append(results, types.FindInTextResults{
						Word:        word,
						MatchedWord: sw,
						LineIndex:   lineNumber,
						CharIndex:   token.CharIndex,
						Context:     context,
					})
				}

			}
		}
	}

	if err := scanner.Err(); err != nil {
		return types.FindInTextResponse{}, err
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].LineIndex != results[j].LineIndex {
			return results[i].LineIndex < results[j].LineIndex
		}
		return results[i].CharIndex < results[j].CharIndex
	})

	response := types.FindInTextResponse{
		Results:       results,
		CaseSensitive: options.CaseSensitive,
		WholeWordOnly: options.WholeWordOnly,
	}

	return response, nil
}

// LetterFrequency читает текстовый файл целиком и возвращает частоту
// каждой буквы в тексте, без учёта регистра. Небуквенные символы
// (пунктуация, цифры, пробелы, переносы строк) в подсчёт не входят.
//
// Результат отсортирован по убыванию частоты, при равной частоте буквы
// упорядочиваются по алфавиту, чтобы результат был одинаковым при
// каждом вызове.
//
// Возвращает ошибку, если файл не удалось открыть или прочитать.
func LetterFrequency(name string) ([]types.LetterCount, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	text := string(data)

	letters := make(map[rune]int)
	for _, r := range text {
		if !unicode.IsLetter(r) {
			continue
		}
		letters[unicode.ToLower(r)]++
	}

	results := make([]types.LetterCount, 0, len(letters))
	for letter, count := range letters {
		results = append(results, types.LetterCount{Letter: letter, Count: count})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Letter < results[j].Letter
	})

	return results, nil
}