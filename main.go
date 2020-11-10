package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sobczak-m/time/rate"
)

func req(l *rate.Limiter, i int) error {
	fmt.Printf("-----------------------------  request: %d\n", i)
	if !l.Allow() {
		return errors.New("rate limit exceeded")
	}
	fmt.Println("-----------------------------------------------")
	return nil
}

func run(l *rate.Limiter, t *time.Ticker, count int) error {
	idx := 0
	for range t.C {
		idx++
		err := req(l, idx)

		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return err
		}
		if idx >= count {
			break
		}

	}
	return nil
}

func main() {
	l := rate.NewLimiter(5, 5)
	t := time.NewTicker(100 * time.Millisecond)
	run(l, t, 20)
}
