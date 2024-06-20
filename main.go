package main

import (
	"github.com/caseymrm/menuet"
)

func main() {
	go taskExecutor()
	menuet.App().Label = "com.github.gonfva.amialive"
	menuet.App().RunApplication()
}
