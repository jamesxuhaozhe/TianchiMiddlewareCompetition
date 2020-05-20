package conf

// default server port is 8080
var serverPort = "8080"

func Init(port string) {
	serverPort = port
}

func GetPort() string {
	return serverPort
}
