package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/time/rate"
	//"github.com/sobczak-m/time/rate"
)

func req(l *rate.Limiter, i int, burst bool) error {
	fmt.Printf("-----------------------------  request: %d | Burst: %v\n", i, burst)
	if !l.Allow() {
		return errors.New("rate limit exceeded")
	}
	fmt.Println("----------------------------------------------- processed ")
	return nil
}

func run(l *rate.Limiter, t *time.Ticker, burstRequestNumber int, requestNumber int) error {
	idx := 1
	for idx <= burstRequestNumber {
		err := req(l, idx, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return err
		}
		if idx >= requestNumber {
			break
		}
		idx++
	}
	time.Sleep(2000)
	for range t.C {
		err := req(l, idx, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return err
		}
		if idx >= requestNumber {
			break
		}
		idx++
	}

	return nil
}

func main() {
	burst := 5
	limit := 0.5
	requests := 20
	ticker := 10 * 199

	l := rate.NewLimiter(rate.Limit(limit), burst)
	t := time.NewTicker(time.Duration(ticker) * time.Millisecond)
	run(l, t, burst, requests)
}
