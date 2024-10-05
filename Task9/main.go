package main

import (
	"GolangPractice/Task9/workerpool"
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	fileName := "example.txt"
	countWorkers := 3
	counter(fileName, countWorkers)
}

// counter - count number of letters many language in file
func counter(fileName string, count int) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	pool := workerpool.Workers(count)
	reader := bufio.NewReader(file)

	var ua, ru, en, rua int

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error: ", err)
			}
			break
		}

		lineCopy := line
		pool.Submit(func() {
			u, r, e, ra := countLetters(lineCopy)
			ua += u
			ru += r
			en += e
			rua += ra
		})
	}
	pool.Shutdown()
	fmt.Printf("English: %d\nUkrainian: %d\nRussian: %d\nKyrilic: %d\n", en, ua, ru, rua)
}

// countLetters - check what language have this letters and add count
func countLetters(line string) (int, int, int, int) {
	var ua, ru, en, rua int

	for _, r := range line {
		if isUa(r) {
			ua++
		} else if isRu(r) {
			ru++
		} else if isEn(r) {
			en++
		} else if isRuUa(r) {
			rua++
		}
	}
	return ua, ru, en, rua
}

func isUa(r rune) bool {
	return r == 'ґ' || r == 'Ґ' || r == 'є' || r == 'Є' || r == 'і' || r == 'І' || r == 'ї' || r == 'Ї'
}
func isRu(r rune) bool {
	return r == 'ё' || r == 'Ё' || r == 'ъ' || r == 'Ъ' || r == 'ы' || r == 'Ы' || r == 'э' || r == 'Э'
}
func isEn(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}
func isRuUa(r rune) bool {
	return (r >= 'А' && r <= 'я') && !isUa(r) && !isRu(r)
}
