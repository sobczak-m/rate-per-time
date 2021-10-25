package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sobczak-m/time/rate"
)

func req(l *rate.Limiter, i int, burst bool, startTime int64) error {
	now := time.Now()
	requestTime := now.Unix()
	requestDuration := requestTime - startTime
	fmt.Printf("-----------------------------  request: %d | Time: %d | Burst: %v\n", i, requestDuration, burst)
	if !l.Allow() {
		return errors.New("rate limit exceeded")
	}
	fmt.Println("----------------------------------------------- processed ")
	return nil
}

func run(l *rate.Limiter, t *time.Ticker, burstRequestNumber int, requestNumber int, ticker time.Duration) error {
	idx := 1
	now := time.Now()
	startTime := now.Unix()
	for idx <= burstRequestNumber {
		err := req(l, idx, true, startTime)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return err
		}
		if idx >= requestNumber {
			break
		}
		idx++
		time.Sleep(ticker * time.Millisecond)
	}

	for range t.C {
		err := req(l, idx, false, startTime)
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
	var laas_limit float32 = 300
	var ticker time.Duration = 6000
	laas_period := 3600
	requests := 600

	burst := calculateBurst(laas_limit)
	limit := calculateLimitPerSecond(laas_limit, laas_period)

	l := rate.NewLimiter(rate.Limit(limit), burst)
	t := time.NewTicker(ticker * time.Millisecond)
	run(l, t, burst, requests, ticker)
}

func calculateBurst(limit float32) int {
	burst := int(limit)
	if burst >= 1 {
		return burst
	}
	return 1
}

func calculateLimitPerSecond(limit float32, period int) float32 {
	return limit / float32(period)
}
