package types

type TopOptions struct {
	Length           int      // 0 - любая длина
	Limit            int      // сколько слов вернуть
	CaseSensitive    bool     // чувствительность к регистру
	ExcludeStopWords bool     // если true, союзы/предлоги/частицы исключаются
	StopWords        []string // свой список. если nil, используется дефолтный
}

type WordCount struct {
	Word  string
	Count int
}

type Stats struct {
	// Знаки пунктуации
	TotalPunctuation int

	// Буквы
	TotalLetters              int
	TotalLettersWithoutSpaces int
	OnlyWordsLetters          int

	// Цифры и числа
	TotalFigures int // Цифры
	TotalNumbers int // Числа

	// Слова
	TotalWords               int
	TotalSyllables           int
	TotalStopWords           int
	StopWordsPercentage      float64
	UniqueWords              int
	AvgWordLengthByLetters   float64
	AvgWordLengthBySyllables float64
	LongestWord              string
	ShortestWord             string

	// Предложения
	TotalSentences    int
	AvgSentenceLength float64 // в словах
	LongestSentence   string
	ShortestSentence  string

	// Лингвистические харктеристики
	TypeTokenRatio   float64 // Лексическое разнообразие
	ReadabilityScore float64 // Индекс читаемости Флеша
}
