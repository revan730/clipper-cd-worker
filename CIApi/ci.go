package CIApi

import (
	"context"

	"github.com/revan730/clipper-cd-worker/log"
	commonTypes "github.com/revan730/clipper-common/types"
	"google.golang.org/grpc"
)

type CIClient struct {
	gClient commonTypes.CIAPIClient
	log     log.Logger
}

func NewClient(address string, logger log.Logger) *CIClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.LogFatal("Couldn't connect to CI gRPC", err)
	}

	c := commonTypes.NewCIAPIClient(conn)
	client := &CIClient{
		gClient: c,
		log:     logger,
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
