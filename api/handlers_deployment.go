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
	}
	err := s.databaseClient.CreateDeployment(deployment)
	if err != nil {
		s.logError("Create deployment error", err)
		return &commonTypes.Empty{}, status.New(http.StatusInternalServerError, "").Err()
	}
	return &commonTypes.Empty{}, nil
	// TODO: Start deployment job
}
