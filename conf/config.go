package conf

// default server port is 8080
var serverPort = "8080"

var datasourcePort = ""

// SetServerPort sets the current port of the server
func SetServerPort(port string) {
	serverPort = port
}

// GetServerPort gets the current port of the server
func GetServerPort() string {
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

// GetLocalTestPort gets
func GetLocalTestPort() string {
	return "80"
}
