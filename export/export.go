package export

import (
	"fmt"
	"encoding/json"
	"os"

	"github.com/fatih/color"
)

var success = color.New(color.FgHiWhite, color.Bold, color.BgBlue)

// ExportJSON сериализует любой результат анализа (types.Stats,
// []types.WordCount, types.NGramResponse и так далее) в JSON-файл по
// указанному пути. Единственная универсальная функция для всех типов
// результатов библиотеки, специфичного кода под каждый тип не требуется,
// поскольку json.MarshalIndent сериализует структуру через рефлексию по
// экспортируемым полям, используя JSON-теги, заданные в types.
//
// Файл записывается с отступами в два пробела для читаемости человеком
// при открытии напрямую, а не одной сплошной строкой.
//
// Параметры:
//   - data: значение любого типа для сериализации, обычно результат
//     одной из функций пакета services (Stats, []WordCount, NGramResponse
//     и подобные).
//   - path: путь к файлу, куда будет записан результат. Если файл уже
//     существует, он будет перезаписан (см. os.WriteFile).
//
// Возвращает ошибку, если data не удалось сериализовать в JSON (например,
// при наличии в структуре несериализуемых полей вроде каналов или
// функций) или если файл не удалось записать.
func ExportJSON(data any, path string) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println()
	msg := fmt.Sprintf("Данные успешно записаны по пути \"%s\"", path)
	success.Print(msg)
	fmt.Println()
	return os.WriteFile(path, bytes, 0644)
}

// ExportTXT сохраняет результат работы Formatting-функции (GetTopFormatting,
// GetStatsFormatting, FindInTextFormatting и так далее) в текстовый файл
// вместо вывода в терминал. Сама Formatting-функция при этом не меняется
// и не знает, что пишет в файл, а не на экран.
//
// Работает через временную подмену os.Stdout на файл на время выполнения
// printFn: поскольку все Formatting-функции пишут через fmt.Println и
// подобные, перенаправление os.Stdout заставляет их вывод попасть в файл
// без единой правки в их коде. На время записи также отключается цвет
// ANSI (color.NoColor = true), поскольку escape-коды цвета делают текст
// нечитаемым при открытии обычным текстовым редактором. И os.Stdout, и
// color.NoColor гарантированно восстанавливаются в исходное состояние
// через defer, даже если printFn запаникует.
//
// Поскольку os.Stdout и color.NoColor это глобальное состояние процесса,
// ExportTXT не потокобезопасна и не должна вызываться параллельно из
// нескольких горутин одновременно, а также не должна вызываться из
// самой printFn повторно (вложенный вызов).
//
// Параметры:
//   - path: путь к текстовому файлу, куда будет записан результат. Если
//     файл уже существует, он будет перезаписан (см. os.Create).
//   - printFn: функция без аргументов, которая внутри себя вызывает один
//     из Formatting с уже готовыми данными, например:
//     func() { formatting.GetTopFormatting(results) }.
//
// Возвращает ошибку, если файл не удалось создать. Ошибки самой printFn
// (например, если Formatting-функция сама может паниковать на
// неожиданных данных) этой функцией не перехватываются.
func ExportTXT(path string, printFn func()) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	oldStdout := os.Stdout
	oldNoColor := color.NoColor
	oldColorOutput := color.Output

	os.Stdout = file
	color.NoColor = true
	color.Output = file

	defer func() {
		os.Stdout = oldStdout
		color.NoColor = oldNoColor
		color.Output = oldColorOutput
	}()

	printFn()
	
	os.Stdout = oldStdout
	color.NoColor = oldNoColor
	color.Output = oldColorOutput

	fmt.Println()
	success.Printf("Данные успешно записаны по пути \"%s\"", path)
	fmt.Println()

	return nil
}