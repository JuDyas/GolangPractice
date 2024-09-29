package main

import "fmt"

func main() {
	A := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	B := [][]int{
		{7, 8},
		{9, 10},
		{11, 12},
	}

	rezult, error := matrix(A, B)
	if error != nil { //Проверка на присутствие ошибки
		fmt.Println(error)
		return
	}
	for _, i := range rezult { //Вывод результата
		fmt.Println(i)
	}
}

func matrix(A, B [][]int) ([][]int, error) {
	var vertA, vertB = len(A), len(B)
	var gorA, gorB = len(A[0]), len(B[0]) //Узнаём количество столбцов/рядов

	if gorA != vertB { //Проверка вохможно ли перемножить матрицы
		return nil, fmt.Errorf("It is impossible to multiply")
	}

	C := make([][]int, vertA) //Создаём матрицу в которую будет записан результат
	for i := range C {
		C[i] = make([]int, gorB)
	}

	for i := 0; i < vertA; i++ { //Проход по рядам A
		for j := 0; j < gorB; j++ { //Проход по строкам B
			sum := 0
			for k := 0; k < gorA; k++ { //Проход по строкам А
				sum += A[i][k] * B[k][j] //Перемножаем эллементы матриц
			}
			C[i][j] = sum //Записываем результат умножения в новую матрицу
		}
	}
	return C, nil
}
