package main

import (
	"time"

	"math/rand"

	"github.com/caseymrm/menuet"
	"github.com/gen2brain/beeep"
	probing "github.com/prometheus-community/pro-bing"
)

const (
	numRuns = 10
)

type Pinger struct {
	Address string
	Pinger  *probing.Pinger
}

func helloPinger() {
	addresses := []string{"8.8.8.8", "1.1.1.1", "9.9.9.9"}
	var sum time.Duration
	var denom int64
	var title string
	lastFewRuns := make([]time.Duration, numRuns)
	currentPointer := 0
	for {
		address := addresses[rand.Intn(len(addresses))]
		pinger, err := probing.NewPinger(address)
		if err != nil {
			panic(err)
		}
		pinger.Count = 1
		pinger.Run()
		stats := pinger.Statistics()
		pinger.Stop()
		currentRtt := stats.AvgRtt

		if denom < numRuns {
			denom += 1
		} else {
			sum = sum - lastFewRuns[currentPointer]
		}
		lastFewRuns[currentPointer] = currentRtt
		sum += currentRtt

		currentPointer += 1
		if currentPointer >= numRuns {
			currentPointer = 0
		}
		avg := time.Duration(int64(sum) / denom)
		if currentRtt > 2*avg || currentRtt == 0 || stats.PacketLoss > 0 {
			title = "**************************************"
			err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
			if err != nil {
				panic(err)
			}
		} else {
			title = avg.String()
		}
		//log.Println(lastFewRuns, currentRtt, sum, avg)
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: title,
		})
		time.Sleep(1 * time.Second)
	}
}

func main() {
	go helloPinger()
	menuet.App().RunApplication()
}
