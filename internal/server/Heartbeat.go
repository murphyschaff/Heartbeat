package Heartbeat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
	probing "github.com/prometheus-community/pro-bing"
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
	//run scheduled jobs
	go scheduler()
	//start API
	api()
	return nil
}

// gathers data for each client based each time the specified scheduler runs
func scheduler() error {
	ctx, cancel := context.WithCancel(context.Background())
	sysLog := NewLogger("Heartbeat System", CONFIGURATION.LogDirectoryPath+"/system.txt", true)
	sysLog.Log("INFO", "Starting up scheduled tasks", false)

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
	logger := NewLogger(client.Name, CONFIGURATION.LogDirectoryPath+"/"+client.Name+".txt", true)

	//check ping on device before attempting SNMP
	ping, err := RunPing(client.IP)
	if !ping || err != nil {
		message := fmt.Sprintf("ping down on client '%s'", client.Name)
		logger.Log("DOWN", message, true)
		return nil
	}

	//ping recieved
	snmp := &gosnmp.GoSNMP{Target: client.IP}
	err = snmp.Connect()
	if err != nil {
		message := fmt.Sprintf("unable to connect to snmp: %s", err)
		logger.Log("DOWN", message, true)
		return err
	}
	defer snmp.Conn.Close()

	//grab data for each monitor for the client
	result, err := snmp.Get(client.GetOIDS())
	if err != nil {
		message := fmt.Sprintf("errors in getting data from snmp client: %s", err)
		logger.Log("ERROR", message, false)
		return err
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

// function that pings IP to see if it is currently available. Returns true if available
func RunPing(IP string) (bool, error) {
	ping, err := probing.NewPinger(IP)

	if err != nil {
		return false, fmt.Errorf("unable to create ping object: %s", err)
	}
	ping.Count = 1
	ping.Timeout = 3 * time.Second
	ping.SetPrivileged(true)

	err = ping.Run()
	if err != nil {
		return false, fmt.Errorf("unable to run ping command: %s", err)
	}

	result := ping.Statistics()
	if result.PacketsRecv > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
