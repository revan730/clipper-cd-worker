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
	changeImageChan chan types.Deployment
	scaleChan       chan types.Deployment
	reInitChan      chan types.Deployment
	deleteChan      chan types.Deployment
}

func NewServer(config Config, logger log.Logger, dbClient db.DatabaseClient) *Server {
	server := &Server{
		config:          config,
		log:             logger,
		databaseClient:  dbClient,
		deploymentsChan: make(chan types.Deployment),
		changeImageChan: make(chan types.Deployment),
		scaleChan:       make(chan types.Deployment),
		reInitChan:      make(chan types.Deployment),
		deleteChan:      make(chan types.Deployment),
	}
	return server
}

// TODO: Too many channels

// GetDepsChan returns read only channel of Deployment type
// used to inform cd worker about new deployments to init
func (s *Server) GetDepsChan() <-chan types.Deployment {
	return s.deploymentsChan
}

// GetImageChangeChan returns read only channel of Deployment type
// used to inform cd worker about image changes in deployments
func (s *Server) GetImageChangeChan() <-chan types.Deployment {
	return s.changeImageChan
}

// GetScaleChan returns read only channel of Deployment type
// used to inform cd worker about deployments to be scaled
func (s *Server) GetScaleChan() <-chan types.Deployment {
	return s.scaleChan
}

// GetReInitChan returns read only channel of Deployment type
// used to inform cd worker about deployments to be reinitialized
// with new manifest
func (s *Server) GetReInitChan() <-chan types.Deployment {
	return s.reInitChan
}

// GetDeleteChan returns read only channel of Deployment type
// used to inform cd worker about deployments to be deleted
func (s *Server) GetDeleteChan() <-chan types.Deployment {
	return s.deleteChan
}

// Run starts api server
func (s *Server) Run() {
	defer s.databaseClient.Close()
	rand.Seed(time.Now().UnixNano())
	err := s.databaseClient.CreateSchema()
	if err != nil {
		s.log.Fatal("Failed to create database schema", err)
	}
	s.log.Info(fmt.Sprintf("Starting api server at port %d", s.config.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		s.log.Fatal("API server failed", err)
	}
	grpcServer := grpc.NewServer()
	commonTypes.RegisterCDAPIServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		s.log.Fatal("failed to serve gRPC API", err)
	}
}
