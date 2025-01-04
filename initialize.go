package main

import (
	"fmt"

	"github.com/murphyschaff/go-helpers"
)

func init(settings *Settings) {
	fmt.Println("Initializing first boot...\nSelect the settings you wish to change")
	//fill out initial settings, or use default settings
	default_message := "Would you like to use the following default settings?\nTime between ping: 5 minutes\n"
	if YN(default_message) {
		settings.SecBtwPct = 300
	} else {
		settings.SecBtwPct = helpers.CorrectStringValidate("Seconds between packets")
	}

}
