package conf

// default server port is 8080
var serverPort = "8080"

var datasourcePort = ""

// SetServerPort sets the current port of the server
func SetSeverPort(port string) {
	serverPort = port
}

// GetPort gets the current port of the server
func GetPort() string {
	return serverPort
}

// SetDatasourcePort sets the port of the data source
func SetDatasourcePort(port string) {
	datasourcePort = port
}

// GetDatasourcePort gets
func GetDatasourcePort() string {
	return datasourcePort
}
