package utils

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