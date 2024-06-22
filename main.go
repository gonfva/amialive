package main

import (
	"github.com/caseymrm/menuet"
)

const (
	TripleAverage = iota
	LessThan250ms
)

func main() {
	go taskExecutor()

	menuet.App().Label = "com.github.gonfva.amialive"
	menuet.App().Children = menuItems
	menuet.App().RunApplication()
}

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}
	items = append(items, menuet.MenuItem{
		Text: "Triple average",
		Clicked: func() {
			menuet.Defaults().SetInteger("AlertOn", TripleAverage)
		},
		State: menuet.Defaults().Integer("AlertOn") == TripleAverage,
	})
	items = append(items, menuet.MenuItem{
		Text: "Less than 250ms",
		Clicked: func() {
			menuet.Defaults().SetInteger("AlertOn", LessThan250ms)
		},
		State: menuet.Defaults().Integer("AlertOn") == LessThan250ms,
	})

	return items
}
