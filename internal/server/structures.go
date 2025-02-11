package Heartbeat

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	CONFIGURATION *Configuration
	DATA          *ClientData
	RUN           = true
)

type Configuration struct {
	//General
	MonitorInterval int    `env:"MONITOR_INTERVAL,required" json:"monitor_interval"` //interval in minutes
	APIPort         string `env:"API_PORT,required" json:"api_port"`
	//Logging
	EmailDomain        string `env:"EMAIL_DOMAIN,required" json:"email_domain"`
	LoggingDestination string `env:"DESTINATION_DOMAIN,required" json:"logging_destination"`
	LogDirectoryPath   string `env:"LOG_DIR_PATH,required" json:"logging_path"`

	//HA
	EnableHA  bool   `env:"HA_ENABLE,required" json:"enable_ha"`
	IsPrimary bool   `env:"IS_PRIMARY" json:"is_primary"`
	HAPeer    string `env:"PEER_IP" json:"ha_peer"`
	HATimer   int    `env:"PEER_TIMER" json:"ha_timer"`

	//Data
	ClientDataFile string `env:"CLIENT_DATA_FILE" json:"client_data_file"`
}

// Heartbeat ClientData: wrapper struct to hold and control current client data
type ClientData struct {
	Clients []Client   `json:"clients"`
	Mutex   sync.Mutex `json:"-"`
}

// reads client data from file
func (d *ClientData) Load() error {
	file, err := os.Open(CONFIGURATION.ClientDataFile)
	if err != nil {
		return fmt.Errorf("unable to open client data file: %s", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&d.Clients)
	if err != nil {
		return fmt.Errorf("unable to parse data from client data file: %s", err)
	}
	return nil
}

// writes client data to data file
func (d *ClientData) Save() error {
	file, err := os.Open(CONFIGURATION.ClientDataFile)
	if err != nil {
		return fmt.Errorf("unable to open client data file: %s", err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return fmt.Errorf("unable to marshal data to file: %s", err)
	}
	file.Write(data)
	return nil
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

// Structure that represents a endpoint to monitor
type Monitor struct {
	Name             string `json:"name"`
	OID              string `json:"oid"`
	CriticalValue    int    `json:"critical_value"`
	CriticalOperator bool   `json:"critical_operator"` //greater than is true
	WarningValue     int    `json:"warning_value"`
	WarningOperator  bool   `json:"warning_operator"` //greater than is true
	AlertWarning     bool   `json:"alert_warning"`
}
