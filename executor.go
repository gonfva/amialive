package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/caseymrm/menuet"
	"github.com/gen2brain/beeep"
	probing "github.com/prometheus-community/pro-bing"
)

const (
	numRuns      = 10
	safetyMargin = 2
	address      = "8.8.8.8"
)

type Stats struct {
	LastNRuns      []int64
	MostRecent     int64
	CurrentSum     int64
	NumIterations  int64
	CurrentPointer int
	avg            int64
	mu             sync.Mutex
	PacketLoss     int
}

func (stats *Stats) run() {
	pinger, err := probing.NewPinger(address)
	if err != nil {
		panic(err)
	}
	pinger.Count = 1
	pinger.Run()
	pstats := pinger.Statistics()
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.MostRecent = int64(pstats.MaxRtt)

	if stats.NumIterations < numRuns {
		stats.NumIterations += 1
	} else {
		stats.CurrentSum = stats.CurrentSum - stats.LastNRuns[stats.CurrentPointer]
	}
	stats.LastNRuns[stats.CurrentPointer] = stats.MostRecent
	stats.CurrentSum += stats.MostRecent

	stats.CurrentPointer += 1
	if stats.CurrentPointer >= numRuns {
		stats.CurrentPointer = 0
	}
	stats.avg = stats.CurrentSum / stats.NumIterations
	stats.getTitle()
}

func (stats *Stats) getTitle() {
	alert := false
	title := time.Duration(stats.MostRecent).String()
	if stats.MostRecent == 0 || stats.PacketLoss > 0 {
		alert = true
		title = "PACKET LOSS"
	}
	if menuet.Defaults().Integer("AlertOn") == TripleAverage && stats.MostRecent > 3*stats.avg {
		alert = true
		title = fmt.Sprintf("Triple average -> Current %v Average %v", time.Duration(stats.MostRecent).String(), time.Duration(stats.avg).String())
	}
	if menuet.Defaults().Integer("AlertOn") == LessThan250ms && time.Duration(stats.MostRecent) > 250*time.Millisecond {
		alert = true
		title = fmt.Sprintf("More than 250ms -> Current %v", time.Duration(stats.MostRecent).String())
	}
	if alert {

		err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		if err != nil {
			panic(err)
		}
		log.Println("OOOPS", stats)
	}
	log.Println(stats)
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: title,
	})
}

func taskExecutor() {

	stats := Stats{LastNRuns: make([]int64, numRuns), CurrentSum: 0, NumIterations: 0}

	for {

		go stats.run()

		time.Sleep(1 * time.Second)
	}
}
