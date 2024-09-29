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
	matrix(A, B)
}

func matrix(A, B [][]int) {
	var vertA = len(A)
	var gorA, gorB = len(A[0]), len(B[0])
	if vertA == gorB {
		C := make([][]int, vertA)
		for i := range C {
			C[i] = make([]int, gorB)
		}

		for i := 0; i < vertA; i++ {
			for j := 0; j < gorB; j++ {
				sum := 0
				for k := 0; k < gorA; k++ {
					sum += A[i][k] * B[k][j]
				}
				C[i][j] = sum
			}
		}
		fmt.Println(C)
	} else {
		fmt.Println("Невозможно перемножить :((")
	}
}
