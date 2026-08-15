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
	TotalSymbols              int
	TotalSymbolsWithoutSpaces int
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

type FindInTextOptions struct {
	Words         []string
	CaseSensitive bool
	WholeWordOnly bool
}

type FindInTextResults struct {
	Word        string
	MatchedWord string
	LineIndex   int
	CharIndex   int
	Context     string
}

type FindInTextResponse struct {
	Results       []FindInTextResults
	CaseSensitive bool
	WholeWordOnly bool
}

type WordWithPosition struct {
	Word      string
	CharIndex int
}
