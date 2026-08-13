package utils

import (
    "os"
	"fmt"
    "text/tabwriter"

	"github.com/westside-jpg/WordMine/types"
    "github.com/fatih/color"
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
        fmt.Println("Нет данных для отображения")
        return
    }

	fmt.Printf(
		"Топ-%d %s %s %s в тексте\n",
		len(results),
		DeclinationWord(len(results), "самое", "самых", "самых"),
		DeclinationWord(len(results), "встречающееся", "встречающихся", "встречающихся"),
		DeclinationWord(len(results), "слово", "слова", "слов"),
	)

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
			DeclinationWord(res.Count, "раз", "раза", "раз"),
			coloredBar,
		)
	}

	w.Flush()
}