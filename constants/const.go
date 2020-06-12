package constants

const (
	// defines the common url end point prefix.
	CommonUrlPrefix = "localhost:"

	ClientProcessPort1 = "8000"

	ClientProcessPort2 = "8001"

	BackendProcessPort1 = "8002"

	// max amount of ping times during Bootstrap
	PingCount = 5

	// Data specific trait tells us that 20000 is a good choice
	BatchSize = 20000

	// Ideally when two clients both finish, the process count should be 2
	ExpectedProcessCount = 2
)
