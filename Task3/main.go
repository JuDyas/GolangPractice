package main

import "fmt"

func main() {
	sl := []int{3, 7, 9, 1, 34, 12, 8, 4, 33, 28, 0}
	bubbleSort(sl)
}
func bubbleSort(sl []int) {
	check := true
	for check {
		check = false
		for i := range sl {
			if i == 0 {
				continue
			} else if sl[i] < sl[i-1] {
				sl[i-1], sl[i] = sl[i], sl[i-1]
				check = true
			}

		}

	}
	fmt.Println(sl)
}
