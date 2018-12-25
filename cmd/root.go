package cmd

import (
	"fmt"
	"os"

	"github.com/revan730/clipper-cd-worker/log"
	"github.com/revan730/clipper-cd-worker/src"
	"github.com/spf13/cobra"
)

var (
	serverPort int
	rabbitAddr string
	ciAddr     string
	dbAddr     string
	db         string
	dbUser     string
	dbPass     string
	logVerbose bool
)

var rootCmd = &cobra.Command{
	Use:   "clipper-cd",
	Short: "CD worker microservice of Clipper CI\\CD",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start worker",
	Run: func(cmd *cobra.Command, args []string) {
		config := &src.Config{
			Port:          serverPort,
			RabbitAddress: rabbitAddr,
			CIAddress:     ciAddr,
			DBAddr:        dbAddr,
			DB:            db,
			DBUser:        dbUser,
			DBPassword:    dbPass,
			Verbose:       logVerbose,
		}

		logger := log.NewLogger(logVerbose)

		worker := src.NewWorker(config, logger)
		worker.Run()
	},
}

// Execute runs application with provided cli params
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// TODO: Remove short flags
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVarP(&serverPort, "port", "p", 8080,
		"Api gRPC port")
	startCmd.Flags().StringVarP(&rabbitAddr, "rabbitmq", "r",
		"amqp://guest:guest@localhost:5672", "Set redis address")
	startCmd.Flags().StringVarP(&ciAddr, "ci", "g",
		"ci-worker:8080", "Set CI gRPC address")
	startCmd.Flags().StringVarP(&dbAddr, "postgresAddr", "a",
		"postgres:5432", "Set PostsgreSQL address")
	startCmd.Flags().StringVarP(&db, "db", "d",
		"clipper", "Set PostgreSQL database to use")
	startCmd.Flags().StringVarP(&dbUser, "user", "u",
		"clipper", "Set PostgreSQL user to use")
	startCmd.Flags().StringVarP(&dbPass, "pass", "c",
		"clipper", "Set PostgreSQL password to use")
	startCmd.Flags().BoolVarP(&logVerbose, "verbose", "v",
		false, "Show debug level logs",
	)
}
