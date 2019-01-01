package src

import (
	"github.com/revan730/clipper-cd-worker/api"
	"github.com/revan730/clipper-cd-worker/CIApi"
	"github.com/revan730/clipper-cd-worker/db"
	"github.com/revan730/clipper-cd-worker/kubectl"
	"github.com/revan730/clipper-cd-worker/queue"
	"github.com/revan730/clipper-cd-worker/types"
	"github.com/revan730/clipper-cd-worker/log"
)

// Worker holds CI worker logic
type Worker struct {
	config         *Config
	jobsQueue      queue.Queue
	kubectl *kubectl.Kubectl
	databaseClient db.DatabaseClient
	ciClient *CIApi.CIClient
	apiServer      *api.Server
	log         log.Logger
}

// NewWorker creates new copy of worker with provided
// config and rabbitmq client
func NewWorker(config *Config, logger log.Logger) *Worker {
	worker := &Worker{
		config: config,
		log: logger,
	}
	dbConfig := types.PGClientConfig{
		DBUser:     config.DBUser,
		DBAddr:     config.DBAddr,
		DBPassword: config.DBPassword,
		DB:         config.DB,
	}
	dbClient := db.NewPGClient(dbConfig)
	worker.jobsQueue = queue.NewRMQQueue(config.RabbitAddress)
	worker.databaseClient = dbClient
	apiConfig := api.Config{
		Port: config.Port,
	}
	apiServer := api.NewServer(apiConfig, logger, dbClient)
	worker.apiServer = apiServer
	ciClient := CIApi.NewClient(config.CIAddress, logger)
	worker.ciClient = ciClient
	// TODO: kubectl config path from app config
	kubectl := kubectl.NewKCtl("")
	worker.kubectl = kubectl
	return worker
}

// Run starts CD worker
func (w *Worker) Run() {
	go w.apiServer.Run()
	w.startConsuming()
}
