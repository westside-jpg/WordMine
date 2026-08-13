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