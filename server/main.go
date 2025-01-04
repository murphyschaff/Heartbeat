package main

import (
	"fmt"
	"os"
)

const (
	CONFIG_FILE string = "./configuration.json"
)

func main() {
	//start Heatbeat server instance
	_, err := os.Stat(CONFIG_FILE)
	if err != nil {
		err = fmt.Errorf("unable to find configuration file")
	} else {
		err = OpenCONFIGURATION()
		if err != nil {
			err = fmt.Errorf("unable to start Heartbeat: %v", err)
		} else {
			//runs normal program
			err = Heartbeat()
			if err != nil {
				err = fmt.Errorf("error while running Heartbeat: %v", err)
			}
		}
	}

	//declaring program end reason
	if err != nil {
		fmt.Printf("Heartbeat has hit a catorstophic error: %v", err)
		os.Exit(1)
	} else {
		fmt.Println("Heartbeat has closed normally")
	}
}
