package main

import (
	"time"

	"math/rand"

	"github.com/caseymrm/menuet"
	"github.com/go-ping/ping"
)

type Pinger struct {
	Address string
	Pinger  *ping.Pinger
}

func getPingers() (pingers []Pinger) {
	addresses := []string{"8.8.8.8", "1.1.1.1", "9.9.9.9"}
	for _, address := range addresses {
		pinger, err := ping.NewPinger(address)
		if err != nil {
			panic(err)
		}
		pingers = append(pingers, Pinger{Address: address, Pinger: pinger})
	}
	return
}

func getPinger(pingers []Pinger) *ping.Pinger {
	return pingers[rand.Intn(len(pingers))].Pinger
}

func helloClock() {
	pingers := getPingers()
	for {
		pinger := getPinger(pingers)
		pinger.Count = 1
		pinger.Run()
		stats := pinger.Statistics()
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: stats.AvgRtt.String(),
		})
		time.Sleep(1 * time.Second)
	}
}

func main() {
	go helloClock()
	menuet.App().RunApplication()
}
