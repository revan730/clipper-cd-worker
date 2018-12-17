package src

// Config represents configuration for application
type Config struct {
	Port int
	// RabbitAddress is used for rabbitmq connection
	RabbitAddress string
	DBAddr        string
	DB            string
	DBUser        string
	DBPassword    string
	Verbose       bool
}
