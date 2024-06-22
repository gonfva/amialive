package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/caseymrm/menuet"
)

const (
	LessThanMaxRTT = iota
	MultipleAverage
)

func main() {
	go taskExecutor()

	menuet.App().Label = "com.github.gonfva.amialive"
	menuet.App().Children = menuItems
	menuet.Defaults().SetInteger("Multiple", 3)
	menuet.Defaults().SetInteger("MaximumRTT", 250)
	menuet.Defaults().SetString("DNSServer", "8.8.8.8")
	menuet.App().RunApplication()
}

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}
	items = append(items, menuet.MenuItem{
		Text: fmt.Sprintf("Multiple of Avg (%v)...", menuet.Defaults().Integer("Multiple")),
		Clicked: func() {
			response := menuet.App().Alert(menuet.Alert{
				MessageText: "Alert when the RTT is higher than a certain multiple of the average last number of pings",
				Inputs:      []string{"Multiple"},
				Buttons:     []string{"Set", "Cancel"},
			})
			if response.Button == 0 && len(response.Inputs) == 1 && response.Inputs[0] != "" {
				if multiple, err := strconv.Atoi(response.Inputs[0]); err == nil {
					menuet.Defaults().SetInteger("AlertOn", MultipleAverage)
					menuet.Defaults().SetInteger("Multiple", multiple)
				}
			}
		},
		State: menuet.Defaults().Integer("AlertOn") == MultipleAverage,
	})
	items = append(items, menuet.MenuItem{
		Text: fmt.Sprintf("Maximum RTT (%v)...", menuet.Defaults().Integer("MaximumRTT")),
		Clicked: func() {
			response := menuet.App().Alert(menuet.Alert{
				MessageText: "Alert when the RTT is higher than a number of mss",
				Inputs:      []string{"Maximum RTT"},
				Buttons:     []string{"Set", "Cancel"},
			})
			if response.Button == 0 && len(response.Inputs) == 1 && response.Inputs[0] != "" {
				if maximum, err := strconv.Atoi(response.Inputs[0]); err == nil {
					menuet.Defaults().SetInteger("AlertOn", LessThanMaxRTT)
					menuet.Defaults().SetInteger("MaximumRTT", maximum)

				}
			}
		},
		State: menuet.Defaults().Integer("AlertOn") == LessThanMaxRTT,
	})
	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})
	items = append(items, menuet.MenuItem{
		Text: fmt.Sprintf("DNS server (%s)...", menuet.Defaults().String("DNSServer")),
		Clicked: func() {
			response := menuet.App().Alert(menuet.Alert{
				MessageText: "Address of the DNS server",
				Inputs:      []string{"DNS server address"},
				Buttons:     []string{"Set", "Cancel"},
			})
			if response.Button == 0 && len(response.Inputs) == 1 && response.Inputs[0] != "" {
				ip := net.ParseIP(response.Inputs[0])
				if ip.To4() != nil {
					menuet.Defaults().SetString("DNSServer", response.Inputs[0])
				}
			}
		},
	})
	return items
}
