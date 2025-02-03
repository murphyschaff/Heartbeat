package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	//create channels for error logging
	monitor = make(chan string)
	status  = make(chan string)
)

func Heartbeat() error {
	//load API
	api()

	//run scheduled jobs
	scheduler()
	return nil
}

// gathers data for each client based each time the specified scheduler runs
func scheduler() error {
	ctx, cancel := context.WithCancel(context.Background())
	sysLog := NewLogger("Heartbeat System", CONFIGURATION.ServerSettings.LogDirectoryPath)

	//scheduler for monitor data collection
	go RunMonitor(ctx, sysLog)
	//scheduler for collecting configuration updates, and total system status
	go RunCheckStatus(ctx, cancel, sysLog)

	switch {
	case monitor != nil:
		return fmt.Errorf("error in running system monitor: %s", monitor)
	case status != nil:
		return fmt.Errorf("error in running system status: %s", status)
	default:
		return nil
	}
}

func RunMonitor(ctx context.Context, logger *Logger) {
	for {
		select {
		case <-ctx.Done():
			logger.Log("INFO", "System monitor is shutting down normally.", false)
			return
		default:
			var tasks sync.WaitGroup
			tasks.Add(len(CONFIGURATION.Clients))

			for _, client := range CONFIGURATION.Clients {
				//find all the related data for each client
			}

			tasks.Wait()

			time.Sleep(time.Duration(CONFIGURATION.ServerSettings.MonitorInterval))
		}
	}
}

func RunCheckStatus(ctx context.Context, cancel context.CancelFunc, logger *Logger) {
	for {
		select {
		case <-ctx.Done():
			logger.Log("INFO", "System status checker is shutting down normally.", false)
			return
		default:
			//try to read configuration file if there has been an update made
			err := CONFIGURATION.Open()
			if err != nil {
				message := "System status checker is unable to read update to configuration file"
				logger.Log("FATAL", message, true)
				status <- message
				cancel()
			}
			//check to see if shutdown requested
			if !CONFIGURATION.Run {
				logger.Log("INFO", "System status has read shutdown request. Initializing shutdown", false)
				cancel()
				return
			}
		}
	}
}
