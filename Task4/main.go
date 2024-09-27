package main

import "fmt"

func main() {
	calculate()
}

func calculate() {
	var num1, num2, answer float64
	var operator string

	fmt.Println("Enter an expression, for example 2 + 2")
	fmt.Scanf("%f %s %f", &num1, &operator, &num2) //Запрашивает пример

	switch operator { //Выбираем оператор (+, -, * или /)
	case "+":
		answer = num1 + num2
	case "-":
		answer = num1 - num2
	case "*":
		answer = num1 * num2
	case "/":
		if num2 == 0 { //При попытке делить на 0 выводим ошибку
			fmt.Println("Error: You can't divide by 0")
		} else {
			answer = num1 / num2
		}
	default: //Если ввод не соотвествует стандарту выдаём ошибку.
		fmt.Println("Error: Invalid input. Please use +, -, *, or /, and separate by spaces.")
		return
	}
	fmt.Println("Answer:", answer)
}
