package Heartbeat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func api() {
	router := gin.Default()
	grp := router.Group("/api")
	//GET commands
	grp.GET("/config", getConfig)
	grp.GET("/clients", getClients)
	grp.GET("/client/:name", getClientByName)
	grp.GET("/status", getStatus)

	//POST commands
	grp.POST("/config/update", setConfig)
	grp.POST("/clients/new", addClient)
	grp.POST("/client/:name/update", setUpdateClient)
	grp.POST("/shutdown", setShutdown)

	router.Run(CONFIGURATION.APIPort)
}

// Settings
func getConfig(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, CONFIGURATION)
}

func setConfig(c *gin.Context) {
	var newConfig Configuration
	err := c.BindJSON(&newConfig)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "settings could not be updated"})
		return
	}
	CONFIGURATION = &newConfig
	//update env file

	c.IndentedJSON(http.StatusOK, CONFIGURATION)
}

// Clients
func getClients(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, DATA.Clients)
}

func addClient(c *gin.Context) {
	var newClient Client

	err := c.BindJSON(&newClient)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "client could not be added"})
		return
	}
	DATA.Clients = append(DATA.Clients, newClient)
	DATA.Save()
	c.IndentedJSON(http.StatusCreated, DATA.Clients)
}

func getClientByName(c *gin.Context) {
	name := c.Param("name")

	for _, client := range DATA.Clients {
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

	for i, client := range DATA.Clients {
		if client.Name == name {
			DATA.Clients[i] = newClient
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
	RUN = false
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Shutdown initiated"})
}
