package formatting

import (
    "os"
	"unicode/utf8"
	"fmt"
	"strings"
    "text/tabwriter"

	"github.com/westside-jpg/WordMine/utils"
	"github.com/westside-jpg/WordMine/types"
    "github.com/fatih/color"
)

// printError печатает сообщение об ошибке в стандартный вывод, выделенное
// белым текстом на красном фоне, используется единообразно во всех
// функциях formatting при отсутствии данных для отображения.
func printError(msg string) {
	color.New(color.FgHiWhite, color.BgRed, color.Bold).Printf(" %s \n", msg)
}

// wrapWidth — ширина переноса строки для длинных предложений в символах.
const wrapWidth = 100

// wrapText разбивает строку на список строк шириной не более width рун,
// перенос идёт по границам слов (strings.Fields), само слово никогда
// не разрывается посередине, даже если оно длиннее width.
func wrapText(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	var line []string
	lineLen := 0

	for _, word := range words {
		wordLen := utf8.RuneCountInString(word)
		extra := wordLen
		if len(line) > 0 {
			extra++ // пробел перед словом
		}

		if len(line) > 0 && lineLen+extra > width {
			lines = append(lines, strings.Join(line, " "))
			line = nil
			lineLen = 0
			extra = wordLen
		}

		line = append(line, word)
		lineLen += extra
	}

	if len(line) > 0 {
		lines = append(lines, strings.Join(line, " "))
	}

	return lines
}

// printWrapped печатает текст с переносом по словам на ширину wrapWidth,
// каждая строка дополнительно сдвинута отступом indent пробелов, включая
// первую строку.
func printWrapped(s string, indent int) {
	prefix := strings.Repeat(" ", indent)
	for _, line := range wrapText(s, wrapWidth) {
		fmt.Println(prefix + line)
	}
}

// readabilityLevel переводит числовое значение индекса читаемости
// Флеша-Оборневой в текстовое описание уровня сложности текста.
//
// Официальной шкалы интерпретации, адаптированной именно под русский
// язык, не существует, большинство источников используют тот же
// диапазон 0-100, что и оригинальная шкала Флеша, привязанный к
// уровням школьного и вузовского образования. Значения за пределами
// 0-100 возможны и означают текст, ещё более простой либо ещё более
// сложный, чем крайние точки шкалы.
func readabilityLevel(score float64) string {
	switch {
	case score >= 90:
		return "очень легко, уровень 5 класса"
	case score >= 80:
		return "легко, уровень 6 класса"
	case score >= 70:
		return "довольно легко, уровень 7 класса"
	case score >= 60:
		return "средне, уровень 8-9 класса"
	case score >= 50:
		return "средне-сложно, уровень 10-11 класса"
	case score >= 30:
		return "сложно, уровень университета"
	default:
		return "очень сложно, уровень выпускника вуза"
	}
}

// ttrLevel даёт грубую качественную оценку лексического разнообразия
// по значению TypeTokenRatio. Чем выше TTR, тем реже повторяются слова
// в тексте, само по себе значение сильно зависит от длины текста,
// поэтому оценка приблизительная, а не строгая шкала.
func ttrLevel(ttr float64) string {
	switch {
	case ttr >= 50:
		return "высокое, словарь заметно разнообразен"
	case ttr >= 25:
		return "среднее"
	default:
		return "низкое, много повторов слов"
	}
}

// GetTopFormatting печатает в стандартный вывод отформатированную таблицу
// топ-слов с порядковыми номерами, количеством вхождений и цветной шкалой
// относительной частоты в процентах от самого частого слова в списке.
//
// Колонки (номер, слово, количество) выравниваются автоматически через
// text/tabwriter по самому широкому значению в каждой колонке. Шкала
// красится в зависимости от процента: зелёный выше 80%, жёлтый выше 40%,
// красный для остального.
//
// Параметры:
//   - results: список слов с их частотой, ожидается отсортированный по
//     убыванию Count (например, результат services.GetTop). Первый элемент
//     используется как эталон 100% для расчёта шкалы остальных строк.
//     Count каждого элемента всегда положителен, поскольку слово попадает
//     в список только при реальных вхождениях в тексте, поэтому деление
//     на ноль при расчёте процента невозможно.
//
// Если results пуст, выводит сообщение об отсутствии данных и завершает
// работу без ошибки.
//
// Цвет автоматически отключается библиотекой fatih/color, если вывод
// перенаправлен не в интерактивный терминал (например, в файл).
func GetTopFormatting(results []types.WordCount) {
    if len(results) == 0 {
		printError("Нет данных для отображения")
        return
    }

	header := color.New(color.FgHiWhite, color.Bold, color.BgBlue)

	fmt.Println()
	header.Printf(" Топ-%d %s %s %s в тексте ",
		len(results),
		utils.DeclinationWord(len(results), "самое", "самых", "самых"),
		utils.DeclinationWord(len(results), "встречающееся", "встречающихся", "встречающихся"),
		utils.DeclinationWord(len(results), "слово", "слова", "слов"),
	)
	fmt.Println()
	fmt.Println()
	warning := color.New(color.FgHiBlack, color.Italic)
	warning.Println("Данные приблизительны,\nвозможны незначительные\nпогрешности")
	fmt.Println()

	maxCount := results[0].Count

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for n, res := range results {
		percent := res.Count * 100 / maxCount

		bar := ""
		for i := range 10 {
			if (i+1)*10 <= percent {
				bar += "█"
			} else {
				bar += "▒"
			}
		}

        var barColor *color.Color
		switch {
		case percent > 80:
			barColor = color.New(color.FgGreen)
		case percent > 40:
			barColor = color.New(color.FgYellow)
		default:
			barColor = color.New(color.FgRed)
		}

		coloredBar := barColor.Sprintf("%s %3d%%", bar, percent)

		fmt.Fprintf(w, "%d.\t%s\t%d %s\t%s\n",
			n+1,
			res.Word,
			res.Count,
			utils.DeclinationWord(res.Count, "раз", "раза", "раз"),
			coloredBar,
		)
	}

	w.Flush()
}

// GetStatsFormatting печатает в стандартный вывод отформатированную сводку
// статистики текста, сгруппированную по смысловым разделам: символы и
// пунктуация, слова, предложения, итоговые метрики. Числовые поля внутри
// раздела выровнены через text/tabwriter, заголовок раздела выделен цветом.
//
// LongestSentence и ShortestSentence выводятся вне таблицы отдельным
// блоком с переносом по словам (см. wrapText), поскольку полный текст
// предложения не обрезается и может занимать несколько строк, что
// сломало бы выравнивание колонок в tabwriter.
//
// Индекс читаемости и лексическое разнообразие сопровождаются кратким
// текстовым пояснением уровня, см. readabilityLevel и ttrLevel.
//
// Если results является нулевым значением types.Stats, выводит сообщение
// об отсутствии данных и завершает работу без ошибки.
func GetStatsFormatting(results types.Stats) {
	if results == (types.Stats{}) {
		printError("Нет данных для отображения")
		return
	}

	header := color.New(color.FgHiWhite, color.Bold, color.BgBlue)
	section := color.New(color.FgCyan, color.Bold)
	label := color.New(color.FgWhite, color.Bold)

	fmt.Println()
	header.Printf(" Общая статистика относительно текста ")
	fmt.Println()
	fmt.Println()
	warning := color.New(color.FgHiBlack, color.Italic)
	warning.Println("Данные приблизительны,\nвозможны незначительные\nпогрешности")
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	section.Fprintln(w, "Символы и пунктуация")
	fmt.Fprintf(w, "  Всего символов:\t%d\n", results.TotalSymbols)
	fmt.Fprintf(w, "  Символов без пробелов:\t%d\n", results.TotalSymbolsWithoutSpaces)
	fmt.Fprintf(w, "  Символов в словах:\t%d\n", results.OnlyWordsLetters)
	fmt.Fprintf(w, "  Знаков пунктуации:\t%d\n", results.TotalPunctuation)
	fmt.Fprintf(w, "  Цифр:\t%d\n", results.TotalFigures)
	fmt.Fprintf(w, "  Чисел:\t%d\n", results.TotalNumbers)
	fmt.Fprintln(w)

	section.Fprintln(w, "Слова")
	fmt.Fprintf(w, "  Всего слов:\t%d\n", results.TotalWords)
	fmt.Fprintf(w, "  Уникальных слов:\t%d\n", results.UniqueWords)
	fmt.Fprintf(w, "  Стоп-слов:\t%d\n", results.TotalStopWords)
	fmt.Fprintf(w, "  Доля стоп-слов:\t%.1f%%\n", results.StopWordsPercentage)
	fmt.Fprintf(w, "  Слогов:\t%d\n", results.TotalSyllables)
	fmt.Fprintf(w, "  Средняя длина слова (буквы):\t%.2f\n", results.AvgWordLengthByLetters)
	fmt.Fprintf(w, "  Средняя длина слова (слоги):\t%.2f\n", results.AvgWordLengthBySyllables)
	fmt.Fprintf(w, "  Самое длинное слово:\t%s\n", results.LongestWord)
	fmt.Fprintf(w, "  Самое короткое слово:\t%s\n", results.ShortestWord)
	fmt.Fprintln(w)

	section.Fprintln(w, "Предложения")
	fmt.Fprintf(w, "  Всего предложений:\t%d\n", results.TotalSentences)
	fmt.Fprintf(w, "  Средняя длина (в словах):\t%.2f\n", results.AvgSentenceLength)

	w.Flush()

	fmt.Println()
	label.Println("  Самое длинное предложение:")
	printWrapped(results.LongestSentence, 4)
	fmt.Println()
	label.Println("  Самое короткое предложение:")
	printWrapped(results.ShortestSentence, 4)
	fmt.Println()

	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	section.Fprintln(w, "Метрики")
	fmt.Fprintf(w, "  Лексическое разнообразие (TTR):\t%.1f%%\t%s\n",
		results.TypeTokenRatio, ttrLevel(results.TypeTokenRatio))
	fmt.Fprintf(w, "  Индекс читаемости (Флеш-Оборнева):\t%.2f%%\t%s\n",
		results.ReadabilityScore, readabilityLevel(results.ReadabilityScore))

	w.Flush()

	fmt.Println()
	warning.Println("TTR — доля уникальных слов от общего числа слов, чем выше, тем меньше повторов")
	warning.Println("Индекс читаемости — чем выше значение, тем проще воспринимается текст")
}