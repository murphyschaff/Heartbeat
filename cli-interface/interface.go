package main

import (
	"fmt"
	"os"

	"github.com/murphyschaff/go-cli"
)

type HeartbeatInterface struct {
	*cli.BaseInterface
}

func main() {
	fmt.Println("Starting Heartbeat-Helper")

	basecli, err := cli.NewInterface("Heartbeat", "./heartbeat-commands.json")
	cli := &HeartbeatInterface{BaseInterface: basecli}

	if err != nil {
		err = fmt.Errorf("unable to start Heartbeat-Helper: %v", err)
	} else {
		err = cli.Run()
		if err != nil {
			err = fmt.Errorf("error while running interface: %v", err)
		}
	}

	if err != nil {
		fmt.Printf("Error while running Heartbeat-Helper: %v", err)
		os.Exit(1)
	} else {
		fmt.Printf("Heartbeat-Helper closed normally")
	}
}

func (h *HeartbeatInterface) Query(query []string) error {
	fmt.Println("This is the local instance")
	return nil
}
