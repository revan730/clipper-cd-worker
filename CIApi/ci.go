package CIApi

import (
	"context"
	"fmt"

	commonTypes "github.com/revan730/clipper-common/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type CIClient struct {
	gClient commonTypes.CIAPIClient
	logger  *zap.Logger
}

func NewClient(address string, logger *zap.Logger) *CIClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect to CI gRPC: %s", err))
	}

	c := commonTypes.NewCIAPIClient(conn)
	client := &CIClient{
		gClient: c,
		logger:  logger,
	}
	return client
}

func (c *CIClient) GetBuildArtifact(buildID int64) (*commonTypes.BuildArtifact, error) {
	return c.gClient.GetBuildArtifact(context.Background(),
		&commonTypes.BuildArtifact{BuildID: buildID})
}

func (c *CIClient) GetBuildArtifactByID(artifactID int64) (*commonTypes.BuildArtifact, error) {
	return c.gClient.GetBuildArtifactByID(context.Background(),
		&commonTypes.BuildArtifact{ID: artifactID})
}
