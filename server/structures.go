package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	CONFIGURATION Configuration
)

type Configuration struct {
	ServerSettings *Settings `json:"settings"`
	Run            bool      `json:"run"`

	Clients []Client   `json:"client_list"`
	Mutex   sync.Mutex `json:"-"`
}

// Creates new configuration object, loads data from file
func (c *Configuration) Open() error {
	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		return fmt.Errorf("unable to open configuration file: %s", err)
	}
	defer file.Close()
	data, _ := io.ReadAll(file)

	err = json.Unmarshal(data, &c)
	if err != nil {
		return fmt.Errorf("unable unmarshalling settings file: %s", err)
	}

	return nil
}

// saves current configuration structure to config file
func (c *Configuration) Save() error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	json, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal config: %v", err)
	}

	file, err := os.Create(CONFIG_FILE)
	if err != nil {
		return fmt.Errorf("unable to open path file: %v", err)
	}
	defer file.Close()
	_, err = file.Write(json)
	if err != nil {
		return fmt.Errorf("unable to write data to config file: %v", err)
	}
	return nil
}

// Heartbeat server settings
type Settings struct {
	//General
	MonitorInterval int    `json:"monitor_interval"`
	APIServerPath   string `json:"apipath"`
	//Logging
	EmailDomain        string `json:"email_domain"`
	LoggingDestination string `json:"email_dest"`
	LogDirectoryPath   string `json:"log_directory_path"`

	//HA
	IsPrimary bool   `json:"HA_primary"`
	HAPeer    string `json:"HA_peer"`
	HATimer   int    `json:"HA_timer"`
}

// Heartbeat server: used to define client servers
// Operator: true -> greater than, false -> less than
type Client struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Status bool   `json:"status"`

	Monitors []Monitor `json:"monitors"`
}

// checks and alerts on the monitor
func (c *Client) CheckAndAlertMonitorStatus(OID string, value int, logger *Logger) {
	monitor := c.GetMonitorByOID(OID)

	if monitor != nil {
		switch {
		case (monitor.CriticalValue > value && monitor.CriticalOperator) || (monitor.CriticalValue < value && !monitor.CriticalOperator):
			//warning alert is generated
			message := fmt.Sprintf("Sensor '%s' is in critical on '%s': Current value: %d", monitor.Name, c.Name, value)
			logger.Log("CRITICAL", message, true)
		case (monitor.WarningValue > value && monitor.WarningOperator) || (monitor.WarningValue < value && !monitor.WarningOperator):
			message := fmt.Sprintf("Sensor '%s' is down on '%s': Current value: %d", monitor.Name, c.Name, value)
			if monitor.AlertWarning {
				logger.Log("WARNING", message, true)
			} else {
				logger.Log("WARNING", message, false)
			}
		}
	} else {
		message := fmt.Sprintf("unable to find matching OID %s on client '%s'", OID, c.Name)
		logger.Log("INFO", message, false)
	}
}

// gets a string slice of all the oids on the monitor
func (c *Client) GetOIDS() []string {
	var oids []string
	for _, monitor := range c.Monitors {
		oids = append(oids, monitor.OID)
	}
	return oids
}

// returns an Monitor object with a matching OID, nil if not found
func (c *Client) GetMonitorByOID(OID string) *Monitor {
	for _, monitor := range c.Monitors {
		if monitor.OID == OID {
			return &monitor
		}
	}
	return nil
}

type Monitor struct {
	Name             string `json:"name"`
	OID              string `json:"oid"`
	CriticalValue    int    `json:"critical_value"`
	CriticalOperator bool   `json:"critical_operator"`
	WarningValue     int    `json:"warning_value"`
	WarningOperator  bool   `json:"warning_operator"`
	AlertWarning     bool   `json:"alert_warning"`
}
