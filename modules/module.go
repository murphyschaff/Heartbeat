package heartbeat_module

import (
	"fmt"
)

type Module interface {
	//must implement
	Settings()              //funtion that will update the settings for the given module
	SendLog(message string) //function that will send
}

type PacketModule interface {
	Module
	//Implemented with BasePacketModule
	CreateBytecode(path string) //creates bytecode needed to run client-side scripts

	//must implement
	ClientScript() //script that will run client-side
}

type FunctionModule interface {
	Module

	//must implement
	HostScript() //script that calls
}

/*
Base Module structures ~ should be implemented along with custom modules
*/
type BasePacketModule struct {
	ModulePath string
}

// generates bytecode to be executed on client system for logging purposes
func (b *BasePacketModule) CreateBytecode(path string) {

}

func (b *BasePacketModule) Settings() {
	fmt.Println("Implement me.")
}

func (b *BasePacketModule) SendLog(message string) {
	fmt.Println("Implement me.")
}
