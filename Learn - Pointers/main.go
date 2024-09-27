package main

import "fmt"

func main() {
	//Указатели.
	tryPointer1 := 3
	tryPointer2 := 11
	tryPointerArray := [5]int{1, 2, 3, 4, 5}
	tryAddVallue := 3

	doubleValue(&tryPointer1)                      //6
	swap(&tryPointer1, &tryPointer2)               //11 6
	incrementArray(&tryPointerArray, tryAddVallue) //4 5 6 7 8

	fmt.Println("doubleValue: ", tryPointer1)
	fmt.Println("swap: ", tryPointer1, tryPointer2)
	fmt.Println("incrementArray: ", tryPointerArray)
}

func doubleValue(a *int) { //Принимает указатель и меняет значение переменной.
	*a = *a * 2
}

func swap(a *int, b *int) { //Принимает два указателя и меняет из значения местами.
	*a, *b = *b, *a
}

func incrementArray(a *[5]int, b int) { //Принимает массив через указатель и увеличивает каждый его эллемент
	for i := range a {
		a[i] += b
	}
}

func findMax(a *[5]int) *int { //Принимает указатель на массив целых чисел и возвращает указатель на максимальный элемент этого массива.
	b := 0
	for i := range a {
		if a[i] > a[b] {
			b = i
		}
	}
	return &a[b]
}
