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
type Client struct {
	Name     string    `json:"name"`
	IP       string    `json:"ip"`
	Status   bool      `json:"status"`
	Monitors []Monitor `json:"processes"`
}

// defines a particular SNMP monitor and the alert values
// Operator values (>, <) stored as bool (greater than true)
type Monitor struct {
	Name         string `json:"name"`
	OID          string `json:"oid"`
	AlertWarning bool   `json:"alert_warning"`

	DownValue       int  `json:"down_value"`
	DownOperator    bool `json:"down_op"`
	WarningValue    int  `json:"warning_value"`
	WarningOperator bool `json:"warning_op"`
}

// checks if a monitor has reached the down threashold
func (m *Monitor) CheckDown(value int) bool {
	if (m.DownValue > value && m.DownOperator) || (m.DownValue < value && !m.DownOperator) {
		return true
	}
	return false
}

// checks if a monitor has reached the warning threashold
func (m *Monitor) CheckWarning(value int) bool {
	if (m.WarningValue > value && m.WarningOperator) || (m.WarningValue < value && m.WarningOperator) {
		return true
	}
	return false
}
