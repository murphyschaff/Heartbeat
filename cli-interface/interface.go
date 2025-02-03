package main

import (
	"fmt"
	"os"

	"github.com/murphyschaff/go-cli"
)

type HeartbeatInterface struct {
	*cli.BaseInterface
	FunctionMap map[string]func(...string)
}

var (
	HeartbeatInt = &HeartbeatInterface{FunctionMap: map[string]func(...string){
		"get":    GetFunction,
		"set":    SetFunction,
		"status": HeartbeatStatus,
	}}
)

func main() {
	fmt.Println("Starting Heartbeat-Helper")

	basecli, err := cli.NewInterface("Heartbeat", "./heartbeat-commands.json")
	HeartbeatInt.BaseInterface = basecli

	if err != nil {
		err = fmt.Errorf("unable to start Heartbeat-Helper: %v", err)
	} else {
		err = cli.Run(HeartbeatInt)
		if err != nil {
			err = fmt.Errorf("error while running interface: %v", err)
		}
	}

	if err != nil {
		fmt.Printf("Error while running Heartbeat-Helper: %v", err)
		os.Exit(1)
	}
}

func (h *HeartbeatInterface) Query(query []string) error {
	if function, exists := h.FunctionMap[query[0]]; exists {
		function() // Call the function
	} else {
		fmt.Println("Function not found!")
	}
	return nil
}

//~~~~~~~interface running commands~~~~~~~~

func GetFunction(options ...string) {

}

func SetFunction(options ...string) {

}

//gets the current status of the heartbeat service (checks if the API is currently up)
func HeartbeatStatus(NA ...string) {
	resp, err := http.Get()
}
