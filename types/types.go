package types

import "encoding/json"

type TopWordsOptions struct {
	Length           int      `json:"length"`             // 0 - любая длина
	Limit            int      `json:"limit"`              // сколько слов вернуть
	CaseSensitive    bool     `json:"case_sensitive"`     // чувствительность к регистру
	ExcludeStopWords bool     `json:"exclude_stop_words"` // если true, союзы/предлоги/частицы исключаются
	StopWords        []string `json:"stop_words"`         // свой список. если nil, используется дефолтный
}

type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type Stats struct {
	// Знаки пунктуации
	TotalPunctuation int `json:"total_punctuation"`

	// Буквы
	TotalSymbols              int `json:"total_symbols"`
	TotalSymbolsWithoutSpaces int `json:"total_symbols_without_spaces"`
	OnlyWordsLetters          int `json:"only_words_letters"`

	// Цифры и числа
	TotalFigures int `json:"total_figures"` // Цифры
	TotalNumbers int `json:"total_numbers"` // Числа

	// Слова
	TotalWords               int     `json:"total_words"`
	TotalSyllables            int     `json:"total_syllables"`
	TotalStopWords            int     `json:"total_stop_words"`
	StopWordsPercentage       float64 `json:"stop_words_percentage"`
	UniqueWords               int     `json:"unique_words"`
	AvgWordLengthByLetters    float64 `json:"avg_word_length_by_letters"`
	AvgWordLengthBySyllables  float64 `json:"avg_word_length_by_syllables"`
	LongestWord               string  `json:"longest_word"`
	ShortestWord              string  `json:"shortest_word"`

	// Предложения
	TotalSentences    int     `json:"total_sentences"`
	AvgSentenceLength float64 `json:"avg_sentence_length"` // в словах
	LongestSentence   string  `json:"longest_sentence"`
	ShortestSentence  string  `json:"shortest_sentence"`

	// Лингвистические харктеристики
	TypeTokenRatio   float64 `json:"type_token_ratio"`  // Лексическое разнообразие
	ReadabilityScore float64 `json:"readability_score"` // Индекс читаемости Флеша
}

type FindInTextOptions struct {
	Words         []string `json:"words"`
	CaseSensitive bool     `json:"case_sensitive"`
	WholeWordOnly bool     `json:"whole_word_only"`
}

type FindInTextResults struct {
	Word        string `json:"word"`
	MatchedWord string `json:"matched_word"`
	LineIndex   int    `json:"line_index"`
	CharIndex   int    `json:"char_index"`
	Context     string `json:"context"`
}

type FindInTextResponse struct {
	Results       []FindInTextResults `json:"results"`
	CaseSensitive bool                `json:"case_sensitive"`
	WholeWordOnly bool                `json:"whole_word_only"`
}

type WordWithPosition struct {
	Word      string `json:"word"`
	CharIndex int    `json:"char_index"`
}

type LetterCount struct {
	Letter rune `json:"letter"`
	Count  int  `json:"count"`
}

// MarshalJSON переопределяет стандартную сериализацию LetterCount, чтобы
// Letter выводился в JSON как строка-символ ("а"), а не как числовой код
// руны (1072). Без этого метода encoding/json сериализовал бы Letter как
// обычное число, поскольку rune это alias для int32 и стандартный
// маршалер не различает "число" и "символ" на уровне типов.
func (l LetterCount) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Letter string `json:"letter"`
		Count  int    `json:"count"`
	}{
		Letter: string(l.Letter),
		Count:  l.Count,
	})
}

type TopNGramOptions struct {
	N             int  `json:"n"`
	Limit         int  `json:"limit"`
	CaseSensitive bool `json:"case_sensitive"`
}

type NGramCount struct {
	NGram string `json:"n_gram"`
	Count int    `json:"count"`
}

type NGramResponse struct {
	NGrams []NGramCount `json:"n_grams"`
	N      int          `json:"n"`
}