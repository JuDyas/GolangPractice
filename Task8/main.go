package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	fmt.Println(retry(5, time.Second))

}

// Отправляем запрос, если отрицательный - повторяем спустя n
func retry(attempts int, sleep time.Duration) string {
	for i := 0; i < attempts; i++ {
		err := exampleRequest()
		if err == nil {
			return "успех"
		}
		fmt.Printf("Попытка %d\n", i+1)
		time.Sleep(sleep)
	}
	return "ошибка"

}

// Подобие запроса, если будет число от 4 до 9, то получим ошибку.
func exampleRequest() error {
	if rand.Intn(10) > 3 {
		return errors.New("ошибка")
	}
	return nil
}
