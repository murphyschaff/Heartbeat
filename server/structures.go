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

	Clients []Client   `json:"client_list"`
	Mutex   sync.Mutex `json:"-"`
}

// Creates new configuration object, loads data from file
func NewConfiguration() (*Configuration, error) {

	var c Configuration

	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		return nil, fmt.Errorf("unable to open configuration file: %s", err)
	}
	defer file.Close()
	data, _ := io.ReadAll(file)

	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("unable unmarshalling settings file: %s", err)
	}

	return &c, nil
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
	SecBtwPct     int    `json:"seconds_between_packet"`
	APIServerPath string `json:"apipath"`
	//Logging
	EmailDomain        string `json:"email_domain"`
	LoggingDestination string `json:"email_dest"`

	//HA
	IsPrimary bool   `json:"HA_primary"`
	HAPeer    string `json:"HA_peer"`
	HATimer   int    `json:"HA_timer"`

	//Modules
	ModuleList []Module `json:"module_list"`
}

// Heartbeat Module struct ~ used to define modules
type Module struct {
	Name           string `json:"Name"`
	IsPacketModule bool   `json:"pkt_module"`
}

// Heartbeat server: used to define client servers
type Client struct {
	Name      string    `json:"name"`
	IP        string    `json:"ip"`
	Status    bool      `json:"status"`
	Processes []Process `json:"processes"`

	AlertCPUUsage    float64 `json:"alertcpu"`
	AlertTemperature float64 `json:"alerttemp"`
	AlertMemory      uint64  `json:"alertmem"`
	AlertDisk        uint64  `json:"alertdisk"`
}

type Process struct {
	Name string `json:"name"`

	AlertCPUUsage  float64 `json:"alertcpu"`
	AlertMemory    uint64  `json:"alertmem"`
	MustRun        bool    `json:"mustrun"`
	MustForeground bool    `json:"mustforeground"`
}

// ClientData, data collected from https://github.com/shirou/gopsutil
type ClientData struct {
	CPUusage    float64 `json:"cpuusage"`    //cpu.Percent
	Temperature float64 `json:"temperature"` //sensors.TemperaturesWithContext
	MemUsage    uint64  `json:"memusage"`    //mem.VirtualMemory
	DiskUsage   uint64  `json:"diskusage"`   //disk.Usage
}

type ProcessData struct {
	CPUusage   float64 `json:"cpuusage"`   //process.CPUprecent
	MemUsage   uint64  `json:"memusage"`   //process.mempercent
	Running    bool    `json:"running"`    //process.isrunning
	Background bool    `json:"background"` //process.isbackground
}
