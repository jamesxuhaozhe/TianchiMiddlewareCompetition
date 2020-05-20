package utils

import (
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/common/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/common/constants"
)

func IsClientProcess() bool {
	serverPort := conf.GetPort()
	if serverPort == constants.ClientProcessPort1 || serverPort == constants.ClientProcessPort2 {
		return true
	}
	return false
}

func IsBackendProcess() bool {
	serverPort := conf.GetPort()
	if serverPort == constants.BackendProcessPort1 {
		return true
	}
	return false
}