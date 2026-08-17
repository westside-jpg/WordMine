package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/westside-jpg/WordMine/export"
	"github.com/westside-jpg/WordMine/formatting"
	"github.com/westside-jpg/WordMine/services"
	"github.com/westside-jpg/WordMine/types"
)

var errorMsg = color.New(color.BgRed, color.FgWhite)

func printErr(msg string) {
	fmt.Println()
	errorMsg.Print(msg)
	fmt.Println()
}

func handleExport(path string, data any, print func()) {
	if path == "" {
		print()
		return
	}

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		if err := export.ExportJSON(data, path); err != nil {
			printErr(fmt.Sprintf("Ошибка экспорта: %v", err))
			os.Exit(1)
		}
	case ".txt":
		if err := export.ExportTXT(path, print); err != nil {
			printErr(fmt.Sprintf("Ошибка экспорта: %v", err))
			os.Exit(1)
		}
	default:
		printErr("Неподдерживаемый формат экспорта, используй .txt или .json")
		os.Exit(1)
	}
}

func printHelp() {
	header := color.New(color.FgWhite, color.Bold, color.BgBlue)
	section := color.New(color.FgHiCyan, color.Bold)
	command := color.New(color.FgHiWhite, color.Bold)

	fmt.Println()
	header.Print(" WordMine. CLI-инструмент для анализа текста ")
	fmt.Println()
	fmt.Println()

	section.Println("Использование:")
	fmt.Println("  wordmine <команда> [флаги]")
	fmt.Println()

	section.Println("Команды:")

	fmt.Printf("  %s      Найти самые частотные слова\n",
		command.Sprint("top"))
	fmt.Printf("  %s    Получить статистику текста\n",
		command.Sprint("stats"))
	fmt.Printf("  %s     Найти слова в тексте\n",
		command.Sprint("find"))
	fmt.Printf("  %s  Частотность букв\n",
		command.Sprint("letters"))
	fmt.Printf("  %s   Найти наиболее частые n-граммы\n",
		command.Sprint("ngrams"))

	fmt.Println()
	section.Println("Для подробной информации о команде:")
	fmt.Println("  wordmine <команда> -h")
	fmt.Println()

}

func setCommandHelp(fs *flag.FlagSet, commandName string, description string) {
	header := color.New(color.FgHiWhite, color.Bold, color.BgBlue)
	section := color.New(color.FgHiCyan, color.Bold)

	fs.Usage = func() {
		fmt.Fprintln(fs.Output())

		header.Fprintf(
			fs.Output(),
			" WordMine / %s ",
			commandName,
		)

		fmt.Fprintln(fs.Output())
		fmt.Fprintln(fs.Output())

		fmt.Fprintln(fs.Output(), description)
		fmt.Fprintln(fs.Output())

		section.Fprintln(fs.Output(), "Использование:")
		fmt.Fprintf(fs.Output(), "  wordmine %s [флаги]\n", commandName)
		fmt.Fprintln(fs.Output())

		section.Fprintln(fs.Output(), "Флаги:")
		fs.PrintDefaults()

		fmt.Fprintln(fs.Output())
	}
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	if command == "-h" || command == "-help" {
		printHelp()
		return
	}

	switch command {
	case "top":
		runTop(args)
	case "stats":
		runStats(args)
	case "find":
		runFind(args)
	case "letters":
		runLetters(args)
	case "ngrams":
		runNGrams(args)
	default:
		msg := fmt.Sprintf("Неизвестная команда: %v", command)
		printErr(msg)
		os.Exit(1)
	}
}

func runTop(args []string) {
	fs := flag.NewFlagSet("top", flag.ExitOnError)
	file := fs.String("file", "", "Путь к текстовому файлу")
	length := fs.Int("length", 0, "Длина слов для возврата")
	limit := fs.Int("limit", 10, "Сколько слов вернуть")
	caseSensitive := fs.Bool("case-sensitive", false, "Учитывать регистр")
	excludeStop := fs.Bool("exclude-stopwords", false, "Исключить стоп-слова")
	exportPath := fs.String("export", "", "Путь для сохранения в TXT или JSON вместо вывода в терминал")

	setCommandHelp(
		fs,
		"top",
		"Показывает наиболее часто встречающиеся слова в тексте",
	)

	fs.Parse(args)

	if *file == "" {
		printErr("Нужно указать файл (-file)")
		os.Exit(1)
	}

	results, err := services.TopWords(*file, types.TopWordsOptions{
		Length:           *length,
		Limit:            *limit,
		CaseSensitive:    *caseSensitive,
		ExcludeStopWords: *excludeStop,
	})
	if err != nil {
		msg := fmt.Sprintf("Ошибка: %v", err)
		printErr(msg)
		os.Exit(1)
	}

	handleExport(*exportPath, results, func() { formatting.TopWordsFormatting(results) })
}

func runStats(args []string) {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	file := fs.String("file", "", "Путь к текстовому файлу")
	exportPath := fs.String("export", "", "Путь для сохранения в TXT или JSON вместо вывода в терминал")

	setCommandHelp(
		fs,
		"stats",
		"Показывает подробную статистику анализируемого текста",
	)

	fs.Parse(args)

	if *file == "" {
		printErr("Нужно указать файл (-file)")
		os.Exit(1)
	}

	results, err := services.Stats(*file)
	if err != nil {
		msg := fmt.Sprintf("Ошибка: %v", err)
		printErr(msg)
		os.Exit(1)
	}

	handleExport(*exportPath, results, func() { formatting.StatsFormatting(results) })
}

func runFind(args []string) {
	fs := flag.NewFlagSet("find", flag.ExitOnError)
	file := fs.String("file", "", "Путь к текстовому файлу")
	search := fs.String("search", "", "Слова для поиска через запятую")
	caseSensitive := fs.Bool("case-sensitive", false, "Учитывать регистр")
	wholeWordOnly := fs.Bool("whole-word", false, "Целое слово или подстрока")
	exportPath := fs.String("export", "", "Путь для сохранения в TXT или JSON вместо вывода в терминал")

	setCommandHelp(
		fs,
		"find",
		"Ищет указанные слова или подстроки в тексте",
	)

	fs.Parse(args)

	if *file == "" {
		printErr("Нужно указать файл (-file)")
		os.Exit(1)
	}

	if *search == "" {
		printErr("Нужно указать слова для поиска (-search)")
		os.Exit(1)
	}

	words := strings.Split(*search, ",")

	results, err := services.FindInText(*file, types.FindInTextOptions{
		Words:         words,
		CaseSensitive: *caseSensitive,
		WholeWordOnly: *wholeWordOnly,
	})
	if err != nil {
		msg := fmt.Sprintf("Ошибка: %v", err)
		printErr(msg)
		os.Exit(1)
	}

	handleExport(*exportPath, results, func() { formatting.FindInTextFormatting(results) })
}

func runLetters(args []string) {
	fs := flag.NewFlagSet("letters", flag.ExitOnError)
	file := fs.String("file", "", "Путь к текстовому файлу")
	exportPath := fs.String("export", "", "Путь для сохранения в TXT или JSON вместо вывода в терминал")

	setCommandHelp(
		fs,
		"letters",
		"Показывает частотность букв в тексте",
	)

	fs.Parse(args)

	if *file == "" {
		printErr("Нужно указать файл (-file)")
		os.Exit(1)
	}

	results, err := services.LetterFrequency(*file)
	if err != nil {
		msg := fmt.Sprintf("Ошибка: %v", err)
		printErr(msg)
		os.Exit(1)
	}

	handleExport(*exportPath, results, func() { formatting.LetterFrequencyFormatting(results) })
}

func runNGrams(args []string) {
	fs := flag.NewFlagSet("ngrams", flag.ExitOnError)
	n := fs.Int("n", 2, "Количество слов в n-граме")
	limit := fs.Int("limit", 10, "Сколько n-грам вернуть")
	caseSensitive := fs.Bool("case-sensitive", false, "Учитывать регистр")
	file := fs.String("file", "", "Путь к текстовому файлу")
	exportPath := fs.String("export", "", "Путь для сохранения в TXT или JSON вместо вывода в терминал")

	setCommandHelp(
		fs,
		"ngrams",
		"Показывает наиболее часто встречающиеся последовательности слов",
	)

	fs.Parse(args)

	if *file == "" {
		printErr("Нужно указать файл (-file)")
		os.Exit(1)
	}

	results, err := services.TopNGrams(*file, types.TopNGramOptions{
		N:             *n,
		Limit:         *limit,
		CaseSensitive: *caseSensitive,
	})
	if err != nil {
		msg := fmt.Sprintf("Ошибка: %v", err)
		printErr(msg)
		os.Exit(1)
	}

	handleExport(*exportPath, results, func() { formatting.NGramsFormatting(results) })
}
