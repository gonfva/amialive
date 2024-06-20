package main

import (
	"fmt"
	"log"
	"time"

	"math/rand"

	"github.com/caseymrm/menuet"
	"github.com/gen2brain/beeep"
	probing "github.com/prometheus-community/pro-bing"
)

const (
	historicRuns = 10
	safetyMargin = 2
)

func taskExecutor() {
	addresses := []string{"8.8.8.8", "1.1.1.1", "9.9.9.9"}
	var sum time.Duration
	var denom int64
	var title string
	lastFewRuns := make([]time.Duration, historicRuns)
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

		if denom < historicRuns {
			denom += 1
		} else {
			sum = sum - lastFewRuns[currentPointer]
		}
		lastFewRuns[currentPointer] = currentRtt
		sum += currentRtt

		currentPointer += 1
		if currentPointer >= historicRuns {
			currentPointer = 0
		}
		avg := time.Duration(int64(sum) / denom)
		if currentRtt > 2*avg || currentRtt == 0 || stats.PacketLoss > 0 {
			if currentRtt > safetyMargin*avg {
				title = fmt.Sprintf("High RTT -> Current %v Average %v", currentRtt.String(), avg.String())
			} else {
				title = "PACKET LOSS"
			}
			err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
			log.Println("OOOPS", lastFewRuns, currentRtt, sum, avg)
			if err != nil {
				panic(err)
			}
		} else {
			title = currentRtt.String()
		}
		log.Println(lastFewRuns, currentRtt, sum, avg)
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: title,
		})
		time.Sleep(1 * time.Second)
	}
}

func main() {
	go taskExecutor()
	menuet.App().Label = "com.github.gonfva.amialive"
	menuet.App().RunApplication()
}
