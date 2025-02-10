package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	Heartbeat "github.com/murphyschaff/Heartbeat/internal/server"
)

const (
	CONFIG_FILE string = "./configuration.json"
)

func main() {
	//start Heatbeat server instance
	var config Heartbeat.Configuration
	err := godotenv.Load(".env")
	if err != nil {
		err = fmt.Errorf("unable to parse .env file: %s", err)
	} else {
		err = env.Parse(&config)
		if err != nil {
			err = fmt.Errorf("unable to read .env file and start Heartbeat: %s", err)
		} else {
			err = Heartbeat.Heartbeat(&config)
			if err != nil {
				err = fmt.Errorf("error while running Heartbeat: %s", err)
			}
		}
	}

	//declaring program end reason
	if err != nil {
		fmt.Printf("Heartbeat has hit a catorstophic error: %s", err)
		os.Exit(1)
	} else {
		fmt.Println("Heartbeat has closed normally")
	}
}
