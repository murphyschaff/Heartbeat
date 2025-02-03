package main

import (
	"fmt"
	"regexp"

	"github.com/murphyschaff/go-helpers"
)

func init() {

	var settings *Settings
	fmt.Println("Initializing first boot...\nSelect the settings you wish to change")
	//fill out initial settings, or use default settings
	default_message := "Would you like to use the following default settings?\nPing Interval: 5 minutes\nAPI Path: http://127.0.0.1/\n"
	if helpers.YesNo(default_message) {
		settings.SecBtwPct = 300
		settings.APIServerPath = "http://127.0.0.1/"
	} else {
		settings.SecBtwPct = helpers.GetInt("Ping Interval")
		settings.APIServerPath = helpers.CorrectStringValidate("API Path")
	}
	//settings that do not have a default, must always be set
	email, _ := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	settings.EmailDomain = helpers.CorrectStringValidate("Email Domain", email)
	settings.LoggingDestination = helpers.CorrectStringValidate("Logging destination email address", email)

	//set HA settings (optional)
	if helpers.YesNo("Would you like to setup in High Availability? (HA)") {
		ip, _ := regexp.Compile("^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$")
		settings.HAPeer = helpers.CorrectStringValidate("IP address for HA pair", ip)
		settings.HATimer = helpers.GetInt("Health Check Interval")
		settings.IsPrimary = helpers.YesNo("Would you like this to be the primary server?")
	}

	CONFIGURATION.ServerSettings = settings
}
