package utils

import (
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
)

func IsClientProcess() bool {
	serverPort := conf.GetServerPort()
	if serverPort == constants.ClientProcessPort1 || serverPort == constants.ClientProcessPort2 {
		return true
	}
	return false
}

func IsBackendProcess() bool {
	serverPort := conf.GetServerPort()
	if serverPort == constants.BackendProcessPort1 {
		return true
	}
	return false
}
