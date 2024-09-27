package main

import "fmt"

func main() {
	calculate()
}

func calculate() {
	var num1, num2, answer float64
	var operator string

	fmt.Println("Enter an expression, for example 2 + 2")
	fmt.Scanf("%f %s %f", &num1, &operator, &num2)

	if operator == "+" {
		answer = num1 + num2
	}
	if operator == "-" {
		answer = num1 - num2
	}
	if operator == "*" {
		answer = num1 * num2
	}
	if operator == "/" {
		answer = num1 / num2
	}
	
	fmt.Println("Answer:", answer)
}
