package src

import "time"

// Config represents configuration for application
type Config struct {
	Port int
	// RabbitAddress is used for rabbitmq connection
	RabbitAddress string
	RedisAddress  string
	LockTimeout   time.Duration
	CIAddress     string
	DBAddr        string
	DB            string
	DBUser        string
	DBPassword    string
	Verbose       bool
}
