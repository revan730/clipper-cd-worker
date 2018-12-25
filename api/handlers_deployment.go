package api

import (
	"context"
	"net/http"

	"github.com/revan730/clipper-cd-worker/types"
	commonTypes "github.com/revan730/clipper-common/types"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateDeployment(ctx context.Context, in *commonTypes.Deployment) (*commonTypes.Empty, error) {
	deployment := &types.Deployment{
		Branch:     in.Branch,
		RepoID:     in.RepoID,
		ArtifactID: in.ArtifactID,
		K8SName:    in.K8SName,
		Manifest:   in.Manifest,
		Replicas:   in.Replicas,
	}
	err := s.databaseClient.CreateDeployment(deployment)
	if err != nil {
		s.log.LogError("Create deployment error", err)
		return &commonTypes.Empty{}, status.New(http.StatusInternalServerError, "").Err()
	}
	s.deploymentsChan <- *deployment
	return &commonTypes.Empty{}, nil
}
