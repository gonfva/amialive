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
	address := menuet.Defaults().String("DNSServer")
	pinger, err := probing.NewPinger(address)
	if err != nil {
		panic(err)
	}
	pinger.Count = 1
	pinger.Run()
	pstats := pinger.Statistics()
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.calculateStats(pstats)
	title := stats.getTitle(true)
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: title,
	})
}

func (stats *Stats) calculateStats(pstats *probing.Statistics) {
	stats.MostRecent = int64(pstats.MaxRtt)

	stats.PacketLoss = pstats.PacketsSent - pstats.PacketsRecv

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

func (stats *Stats) String() string {
	durationSlice := make([]string, len(stats.LastNRuns))
	for i, v := range stats.LastNRuns {

		if i == stats.CurrentPointer-1 {
			durationSlice[i] = fmt.Sprintf("%s**", time.Duration(v).String())
		} else {
			durationSlice[i] = time.Duration(v).String()
		}
	}
	str := fmt.Sprintf("LastNRuns: %v Average: %s MostRecent: %v", durationSlice, time.Duration(stats.avg).String(), time.Duration(stats.MostRecent).String())
	return str
}

func (stats *Stats) getTitle(withSound bool) string {
	alert := false
	title := time.Duration(stats.MostRecent).String()
	if stats.MostRecent == 0 || stats.PacketLoss > 0 {
		alert = true
		title = "PACKET LOSS"
	}
	if menuet.Defaults().Integer("AlertOn") == MultipleAverage && stats.MostRecent > int64(menuet.Defaults().Integer("Multiple"))*stats.avg {
		alert = true
		title = fmt.Sprintf("Multiple of average -> Current %v Average %v", time.Duration(stats.MostRecent).String(), time.Duration(stats.avg).String())
	}
	if menuet.Defaults().Integer("AlertOn") == LessThanMaxRTT && time.Duration(stats.MostRecent) > time.Duration(int64(menuet.Defaults().Integer("MaximumRTT")))*time.Millisecond {
		alert = true
		title = fmt.Sprintf("More than %vms -> Current %v", menuet.Defaults().Integer("MaximumRTT"), time.Duration(stats.MostRecent).String())
	}
	if alert {
		if withSound {
			err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
			if err != nil {
				panic(err)
			}
		}
		log.Println("ERROR", stats.String())
	} else {
		log.Println(stats)
	}
	return title
}

func taskExecutor() {

	stats := Stats{LastNRuns: make([]int64, numRuns), CurrentSum: 0, NumIterations: 0}

	for {

		go stats.run()

		time.Sleep(1 * time.Second)
	}
}
