package Heartbeat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
)

func Heartbeat(config *Configuration) error {
	CONFIGURATION = config
	//initialize DATA
	DATA = &ClientData{Clients: []Client{}}
	//get client data from file
	err := DATA.Load()
	if err != nil {
		return fmt.Errorf("unable to load initial data: %s", err)
	}
	//load API
	api()

	//run scheduled jobs
	scheduler()
	return nil
}

// gathers data for each client based each time the specified scheduler runs
func scheduler() error {
	ctx, cancel := context.WithCancel(context.Background())
	sysLog := NewLogger("Heartbeat System", CONFIGURATION.LogDirectoryPath, true)

	//scheduler for monitor data collection
	go RunMonitor(ctx, cancel, sysLog)
	return nil
}

func RunMonitor(ctx context.Context, cancel context.CancelFunc, logger *Logger) {
	for {
		select {
		case <-ctx.Done():
			logger.Log("INFO", "System monitor is shutting down normally.", false)
			return
		default:
			logger.Log("INFO", "Starting new monitor job", false)
			var tasks sync.WaitGroup
			tasks.Add(len(DATA.Clients))

			for _, client := range DATA.Clients {
				//find all the related data for each client
				go GetClientData(&client)
			}

			tasks.Wait()
			ctime := time.Now()
			delta := time.Duration(CONFIGURATION.MonitorInterval) * time.Minute
			//check every 10 seconds for the shutdown command
			for time.Since(ctime) >= delta {
				if !RUN {
					logger.Log("INFO", "System monitor is shutting down normally.", false)
					cancel()
					return
				}
				time.Sleep(10 * time.Second)
			}
		}
	}
}

func GetClientData(client *Client) error {
	//establish logger and SNMP client
	logger := NewLogger(client.Name, CONFIGURATION.LoggingDestination+"/"+client.Name+".txt", false)
	snmp := &gosnmp.GoSNMP{Target: client.IP}
	err := snmp.Connect()
	if err != nil {
		message := fmt.Sprintf("unable to connect to snmp: %s", err)
		logger.Log("FATAL", message, true)
		return fmt.Errorf(message)
	}
	defer snmp.Conn.Close()

	//grab data for each monitor for the client
	result, err := snmp.Get(client.GetOIDS())
	if err != nil {
		message := fmt.Sprintf("errors in getting data from snmp client: %s", err)
		logger.Log("ERROR", message, false)
	}

	var values []int

	for _, variable := range result.Variables {
		values = append(values, variable.Value.(int))
		//verify if this is within the warning/down threashold for the set OID
		client.CheckAndAlertMonitorStatus(variable.Name, values[len(values)-1], logger)
	}

	//~~~~FUTURE~~~ Add options for historical data
	return nil
}
