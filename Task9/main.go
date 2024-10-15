package main

import (
	"GolangPractice/Task9/workerpool"
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"unicode"
)

var (
	inputString  = flag.String("i", "dfsііййjhjf342.fуццаавімdїїїsf.sf[q[wыыыіііe", "input string")
	countWorkers = flag.Int("w", 3, "count of workers")
)

func main() {
	flag.Parse()
	l, err := readInput(*inputString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(countLetters(l, *countWorkers))
}

func readInput(inputString string) ([]string, error) {
	reader := bufio.NewReader(strings.NewReader(inputString))
	var letters []string
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if unicode.IsLetter(r) {
			letters = append(letters, string(r))
		}
	}

	return letters, nil
}

func countLetters(letters []string, countWorkers int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	pool := workerpool.Workers(countWorkers)
	var ua, ru, en, rua int
	var result []int

	if len(letters) != 0 {
		for _, l := range letters {
			letter := l
			pool.Submit(func(ctx2 context.Context) {
				select {
				case <-ctx.Done():
					fmt.Println("Задача отменена:", ctx.Err())
					return
				default:
					switch {
					case isUa([]rune(letter)[0]):
						ua++
					case isRu([]rune(letter)[0]):
						ru++
					case isEn([]rune(letter)[0]):
						en++
					case isRuUa([]rune(letter)[0]):
						rua++
					}
				}
			})
		}
	} else {
		return nil, fmt.Errorf("slice is empty")
	}

	pool.Shutdown()
	result = append(result, ua, ru, en, rua)
	return result, nil
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
