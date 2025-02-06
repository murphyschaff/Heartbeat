package Heartbeat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func api() {
	router := gin.Default()
	//GET commands
	router.GET("/settings", getSettings)
	router.GET("/clients", getClients)
	router.GET("/client/:name", getClientByName)
	router.GET("/status", getStatus)

	//POST commands
	router.POST("/settings/update", setSettings)
	router.POST("/clients/new", addClient)
	router.POST("/client/:name/update", setUpdateClient)
	router.POST("/shutdown", setShutdown)

	router.Run(CONFIGURATION.ServerSettings.APIServerPath)
}

// Settings
func getSettings(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, CONFIGURATION.ServerSettings)
}

func setSettings(c *gin.Context) {
	var newSettings Settings
	err := c.BindJSON(&newSettings)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "settings could not be updated"})
		return
	}
	CONFIGURATION.ServerSettings = &newSettings
	CONFIGURATION.Save()
	c.IndentedJSON(http.StatusOK, CONFIGURATION.ServerSettings)
}

// Clients
func getClients(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, CONFIGURATION.Clients)
}

func addClient(c *gin.Context) {
	var newClient Client

	err := c.BindJSON(&newClient)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "client could not be added"})
		return
	}
	CONFIGURATION.Clients = append(CONFIGURATION.Clients, newClient)
	CONFIGURATION.Save()
	c.IndentedJSON(http.StatusCreated, CONFIGURATION.Clients)
}

func getClientByName(c *gin.Context) {
	name := c.Param("name")

	for _, client := range CONFIGURATION.Clients {
		if client.Name == name {
			c.IndentedJSON(http.StatusOK, client)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "client not found"})
}

func setUpdateClient(c *gin.Context) {
	var newClient Client
	name := c.Param("name")

	err := c.BindJSON(&newClient)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "client could not be updated"})
		return
	}

	for i, client := range CONFIGURATION.Clients {
		if client.Name == name {
			CONFIGURATION.Clients[i] = newClient
			c.IndentedJSON(http.StatusCreated, newClient)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "client not found"})
}

func getStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "running"})
}

func setShutdown(c *gin.Context) {
	CONFIGURATION.Run = false
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Shutdown initiated"})
}
