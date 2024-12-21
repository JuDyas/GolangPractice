package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Enter an expression, for example 2 + 2 * 2")
	read := bufio.NewReader(os.Stdin)      //Запрашивает пример
	readText, _ := read.ReadString('\n')   //Ожидание ввода до нажатия enter(/m)
	readText = strings.TrimSpace(readText) // Должно очистить лишние пробелы если они будут (в начале и конце строки)
	fmt.Println(calculate(readText))

}

func calculate(text string) string { //
	splitText := strings.Fields(text) //Разбиваем строку на отдельные элементы (по пробелам)
	splitText = multyply(splitText)   //Деление и умножение
	answer := addsubst(splitText)     //Сочетание и вычетание
	return answer                     //Вывод ответа
}

func multyply(expression []string) []string { //Умножение с делением идут первыми, поэтому сделал отдельные функции чтобы пройтись сначала по ним.
	var result float64 //Переменная сохраняет результат отдельных операций между 2мя числами.

	i := 0
	for i < len(expression) {
		element := expression[i] //Текущий эллемент

		if element == "*" || element == "/" {
			if i-1 >= 0 && i+1 < len(expression) { //Проверка на существвование индекса, воизбежание выхода за массив
				prev, _ := strconv.ParseFloat(expression[i-1], 64)
				next, _ := strconv.ParseFloat(expression[i+1], 64) //Меняем тип данных для предыдущего и следующего элемента после оператора

				if element == "*" {
					result = prev * next
				} else {
					result = prev / next
				}

				expression[i-1] = fmt.Sprintf("%f", result)              // Заменяем предыдущий элемент на результат вычислеи
				expression = append(expression[:i], expression[i+2:]...) // Удаляем оператор и следующий элемент
				i--                                                      // Корректируем индекс чтобы можно было снова пройтись по оператору
			}
		}
		i++
	}
	return expression //Возвращаем изменённый массив
}
func addsubst(expression []string) string {
	var result float64
	result, _ = strconv.ParseFloat(expression[0], 64) //Меняем тип первого эллемента

	i := 1
	for i < len(expression) {
		element := expression[i]

		if i+1 < len(expression) {
			next, _ := strconv.ParseFloat(expression[i+1], 64)

			if element == "+" {
				result += next
			} else {
				result -= next
			}
		}
		i += 2
	}
	return fmt.Sprintf("%f", result) //Вывод результата
}
