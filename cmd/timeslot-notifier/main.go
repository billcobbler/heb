package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/billcobbler/heb-to-go/pkg/heb-api/v1"
)

type config struct {
	zip           string
	miles         int
	interval      time.Duration
	loopOnSuccess bool
}

func main() {
	zip := os.Getenv("HEB_ZIP")
	if zip == "" {
		fmt.Println("HEB_ZIP is required")
		os.Exit(1)
	}
	miles, err := strconv.Atoi(os.Getenv("HEB_MILES"))
	if err != nil {
		fmt.Println("HEB_MILES is required")
	}
	if miles == 0 {
		fmt.Println("HEB_MILES must be an integer greater than 0")
		os.Exit(1)
	}

	interval, err := time.ParseDuration(os.Getenv("HEB_TIMER"))
	if err != nil {
		fmt.Println("HEB_TIMER is required and must a duration string (ex: 1h)")
		os.Exit(1)
	}
	if interval < time.Minute*1 {
		fmt.Println("HEB_TIMER must be greater than one minute (ex: 1m or 5m)")
		os.Exit(1)
	}

	var loopOnSuccess bool
	if os.Getenv("HEB_CONTINUE_ON_SUCCESS") != "" {
		loopOnSuccess = true
	}

	cfg := config{zip, miles, interval, loopOnSuccess}
	heb := heb.NewClient()
	ctx, stop := context.WithCancel(context.Background())
	go run(ctx, stop, heb, cfg)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	select {
	case <-sigc:
		stop()
		fmt.Println("")
		fmt.Println("exited")
	case <-ctx.Done():
	}

}

func run(ctx context.Context, stop context.CancelFunc, heb *heb.Client, cfg config) {
	ok, err := pullSlots(heb, cfg.zip, cfg.miles)
	if err != nil {
		fmt.Println("failed to pull available slot information: ", err)
		stop()
		return
	}
	if ok {
		if !cfg.loopOnSuccess {
			stop()
			return
		}
	}
	ticker := time.NewTicker(cfg.interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ok, err := pullSlots(heb, cfg.zip, cfg.miles)
			if err != nil {
				fmt.Println("failed to pull available slot information: ", err)
				stop()
				return
			}
			if ok {
				if !cfg.loopOnSuccess {
					stop()
					return
				}
			}
		}
	}
}

func pullSlots(heb *heb.Client, zip string, miles int) (bool, error) {
	stores, err := heb.LocateStores(zip, miles)
	if err != nil {
		return false, err
	}
	if len(stores) == 0 {
		return false, errors.New(fmt.Sprintf("no stores within %d mile(s) of zip %s", miles, zip))
	}

	var slotsFound bool
	for _, s := range stores {
		fmt.Printf("ID: %s, Name: %s, Zip: %s\n", s.ID, s.Name, s.PostalCode)
		fmt.Println("==========================")
		timeslots, err := heb.GetStoreTimeslots(s.ID)
		if err != nil {
			return false, err
		}

		dateCounts := timeSlotCountsByDate(timeslots)
		dates := make([]string, 0, len(dateCounts))
		for k := range dateCounts {
			dates = append(dates, k)
		}

		sort.Strings(dates)
		for _, k := range dates {
			fmt.Printf("%s: %d\n", k, dateCounts[k])
		}
		fmt.Println()

		if len(timeslots) > 0 {
			slotsFound = true
		}
	}

	return slotsFound, nil
}

func timeSlotCountsByDate(slots []heb.Timeslot) map[string]int {
	dateCounts := make(map[string]int)
	for _, ts := range slots {
		dateCounts[ts.Date]++
	}

	return dateCounts
}
