package api

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/revan730/clipper-cd-worker/db"
	"github.com/revan730/clipper-cd-worker/types"
	commonTypes "github.com/revan730/clipper-common/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	Port int
}

type Server struct {
	logger          *zap.Logger
	config          Config
	databaseClient  db.DatabaseClient
	deploymentsChan chan types.Deployment
}

func NewServer(config Config, logger *zap.Logger, dbClient db.DatabaseClient) *Server {
	server := &Server{
		config:          config,
		logger:          logger,
		databaseClient:  dbClient,
		deploymentsChan: make(chan types.Deployment),
	}
	return server
}

func (s *Server) logFatal(msg string, err error) {
	defer s.logger.Sync()
	s.logger.Fatal(msg, zap.Error(err))
}

func (s *Server) logError(msg string, err error) {
	defer s.logger.Sync()
	s.logger.Error(msg, zap.String("packageLevel", "api"), zap.Error(err))
}

func (s *Server) logInfo(msg string) {
	defer s.logger.Sync()
	s.logger.Info("INFO", zap.String("msg", msg), zap.String("packageLevel", "api"))
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
		s.logFatal("Failed to create database schema", err)
	}
	s.logger.Info("Starting api server", zap.Int("port", s.config.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		s.logFatal("API server failed", err)
	}
	grpcServer := grpc.NewServer()
	commonTypes.RegisterCDAPIServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		s.logFatal("failed to serve: %s", err)
	}
}
