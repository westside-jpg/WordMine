package utils

import (
    "regexp"
    "strings"
    "unicode"
	"github.com/westside-jpg/WordMine/types"
)
// DeclinationWord возвращает грамматически правильную форму слова
// в зависимости от числительного n, по правилам русского склонения.
//
// Параметры:
//   - n: число, для которого подбирается форма слова (знак не важен, учитывается только модуль)
//   - one: форма для чисел, оканчивающихся на 1 (кроме 11), например "яблоко"
//   - two: форма для чисел, оканчивающихся на 2-4 (кроме 12-14), например "яблока"
//   - many: форма для остальных случаев, например "яблок"
//
// Пример: DeclinationWord(21, "яблоко", "яблока", "яблок") вернёт "яблоко".
// Пример: DeclinationWord(14, "яблоко", "яблока", "яблок") вернёт "яблок" (исключение 11-14).
func DeclinationWord(n int, one, two, many string) string {
    absN := n
    if absN < 0 {
        absN = -absN
    }

    lastTwoDigits := absN % 100

    if lastTwoDigits >= 11 && lastTwoDigits <= 14 {
        return many
    }

    switch absN % 10 {
    case 1:
        return one
    case 2, 3, 4:
        return two
    default:
        return many
    }
}

// sentencePattern находит предложения целиком, включая завершающий
// знак препинания: последовательность любых символов, кроме .!?, за
// которой следует один или несколько знаков .!? подряд (чтобы многоточие
// "..." попадало в конец предложения одним блоком, а не порождало
// пустые совпадения между точками).
var sentencePattern = regexp.MustCompile(`[^.!?]+[.!?]+`)

// SplitIntoSentences разбивает текст на предложения по знакам
// препинания .!?, сохраняя завершающий знак в каждом предложении, и
// возвращает их списком без пустых строк и лишних пробелов по краям.
//
// Параметры:
//   - text: исходный текст для разбиения, может содержать переносы строк
//
// Возвращает слайс предложений в порядке их появления в тексте, каждое
// предложение оканчивается на исходный .!? Если текст пустой или не
// содержит предложений, вернёт пустой слайс. Текст без завершающего
// знака препинания в самом конце (например, обрывается на середине
// фразы) не попадёт в результат, так как не совпадёт с паттерном.
//
// Ограничение: разбиение основано на пунктуации и не распознаёт
// сокращения (например, "т.д.", "А.С. Пушкин") или десятичные дроби
// (например, "3.14"), такие места могут быть ошибочно приняты за конец
// предложения. Для более точной сегментации потребуется список исключений
// или полноценная NLP-библиотека.
func SplitIntoSentences(text string) []string {
    matches := sentencePattern.FindAllString(text, -1)

	sentences := make([]string, 0, len(matches))
	for _, s := range matches {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			sentences = append(sentences, trimmed)
		}
	}
	return sentences
}


const vowels = "аеёиоуыэюя"

// CountSyllables возвращает приблизительное число слогов в слове,
// считая количество гласных букв. Для русского языка это стандартное
// упрощение, поскольку слог почти всегда строится вокруг одной гласной,
// и не требует полноценного морфологического анализа.
func CountSyllables(word string) int {
	count := 0
	lower := strings.ToLower(word)
	for _, r := range lower {
		if strings.ContainsRune(vowels, r) {
			count++
		}
	}
	// у любого непустого слова хотя бы один слог,
	// на случай аббревиатур или слов без явных гласных
	if count == 0 && lower != "" {
		count = 1
	}
	return count
}

// CleanWord удаляет с обоих краёв строки все символы, кроме букв и цифр,
// оставляя внутренние символы нетронутыми (например, дефис в "кто-то"
// сохраняется). Используется как единый способ очистки слова от знаков
// препинания перед подсчётом и сравнением во всех функциях пакета.
//
// Параметры:
//   - word: слово для очистки, может содержать пунктуацию по краям
//
// Возвращает очищенное слово. Если слово целиком состоит из символов,
// не являющихся буквами или цифрами (например, "--" или "..."), вернёт
// пустую строку.
func CleanWord(word string) string {
    word = strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

    return word
}

// SplitLineWithPositions разбивает строку на слова, попутно
// запоминая, на какой позиции начинается каждое слово.
func SplitLineWithPositions(line string) []types.WordWithPosition {
	var result []types.WordWithPosition

	runes := []rune(line)
	i := 0
	for i < len(runes) {
		for i < len(runes) && unicode.IsSpace(runes[i]) {
			i++
		}
		if i >= len(runes) {
			break
		}

		start := i
		for i < len(runes) && !unicode.IsSpace(runes[i]) {
			i++
		}

		result = append(result, types.WordWithPosition{
			Word:      string(runes[start:i]),
			CharIndex: start,
		})
	}

	return result
}