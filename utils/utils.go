package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
)

// IsClientProcess checks if the current process is client.
func IsClientProcess() bool {
	serverPort := conf.GetServerPort()
	if serverPort == constants.ClientProcessPort1 || serverPort == constants.ClientProcessPort2 {
		return true
	}
	return false
}

// IsBackendProcess checks if the current process is backend.
func IsBackendProcess() bool {
	serverPort := conf.GetServerPort()
	if serverPort == constants.BackendProcessPort1 {
		return true
	}
	return false
}

// MD5 returns the md5 hash of the input string.
func MD5(s string) string {
	return md5Bytes([]byte(s))
}

func md5Bytes(s []byte) string {
	ret := md5.Sum(s)
	return hex.EncodeToString(ret[:])
}
