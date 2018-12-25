package api

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/revan730/clipper-cd-worker/db"
	"github.com/revan730/clipper-cd-worker/log"
	"github.com/revan730/clipper-cd-worker/types"
	commonTypes "github.com/revan730/clipper-common/types"
	"google.golang.org/grpc"
)

type Config struct {
	Port int
}

type Server struct {
	log             log.Logger
	config          Config
	databaseClient  db.DatabaseClient
	deploymentsChan chan types.Deployment
}

func NewServer(config Config, logger log.Logger, dbClient db.DatabaseClient) *Server {
	server := &Server{
		config:          config,
		log:             logger,
		databaseClient:  dbClient,
		deploymentsChan: make(chan types.Deployment),
	}
	return server
}

// GetDepsChan returns read only channel of Deployment type
// used to inform cd worker about new deployments to init
func (s *Server) GetDepsChan() <-chan types.Deployment {
	return s.deploymentsChan
}

// Run starts api server
func (s *Server) Run() {
	defer s.databaseClient.Close()
	rand.Seed(time.Now().UnixNano())
	err := s.databaseClient.CreateSchema()
	if err != nil {
		s.log.LogFatal("Failed to create database schema", err)
	}
	s.log.LogInfo(fmt.Sprintf("Starting api server at port %d", s.config.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		s.log.LogFatal("API server failed", err)
	}
	grpcServer := grpc.NewServer()
	commonTypes.RegisterCDAPIServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		s.log.LogFatal("failed to serve gRPC API", err)
	}
}
