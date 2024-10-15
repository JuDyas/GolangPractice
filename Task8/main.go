package main

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

type Request interface {
	Request() error
}

type exampleRequest struct{}

func (r *exampleRequest) Request() error {
	switch rand.Intn(10) {
	case 0, 2, 4:
		return nil
	case 1, 3, 5:
		return errors.New("bad request")
	case 6, 7, 8, 9:
		return errors.New("server starting up")
	}

	return nil
}

func main() {
	rand.NewSource(time.Now().UnixNano())
	err := retry(5, time.Second, &exampleRequest{})
	if err != nil {
		log.SetPrefix("ERROR: ")
		log.Fatal(err)
	} else {
		log.Printf("INFO: OK")
	}
}

// retry - send request many times
func retry(attempts int, sleep time.Duration, r Request) error {
	for i := 0; i < attempts; i++ {
		err := r.Request()
		if err == nil {
			return nil
		} else if err.Error() == "bad request" {
			return err
		} else {
			time.Sleep(sleep)
		}
	}

	return errors.New("all attempts failed")
}
