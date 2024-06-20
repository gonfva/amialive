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
}

func (stats *Stats) getTitle() (title string) {
	if stats.MostRecent > 2*stats.avg || stats.MostRecent == 0 || stats.PacketLoss > 0 {
		if stats.MostRecent > safetyMargin*stats.avg {
			title = fmt.Sprintf("High RTT -> Current %v Average %v", time.Duration(stats.MostRecent).String(), time.Duration(stats.avg).String())
		} else {
			title = "PACKET LOSS"
		}
		err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		log.Println("OOOPS", stats)
		if err != nil {
			panic(err)
		}
	} else {
		title = time.Duration(stats.MostRecent).String()
	}
	log.Println(stats)
	return title
}

func taskExecutor() {

	var title string

	stats := Stats{LastNRuns: make([]int64, numRuns), CurrentSum: 0, NumIterations: 0}

	for {

		go stats.run()
		title = stats.getTitle()

		menuet.App().SetMenuState(&menuet.MenuState{
			Title: title,
		})
		time.Sleep(1 * time.Second)
	}
}
