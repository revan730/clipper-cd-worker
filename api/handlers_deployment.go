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
	// TODO: Handle 'already exists' error
	if err != nil {
		s.log.Error("Create deployment error", err)
		return &commonTypes.Empty{}, status.New(http.StatusInternalServerError, "").Err()
	}
	s.deploymentsChan <- *deployment
	return &commonTypes.Empty{}, nil
}

func (s *Server) ChangeImage(ctx context.Context, in *commonTypes.Deployment) (*commonTypes.Empty, error) {
	deployment := &types.Deployment{
		ID:         in.ID,
		ArtifactID: in.ArtifactID,
	}
	s.changeImageChan <- *deployment
	return &commonTypes.Empty{}, nil
}

func (s *Server) ScaleDeployment(ctx context.Context, in *commonTypes.Deployment) (*commonTypes.Empty, error) {
	deployment := &types.Deployment{
		ID:       in.ID,
		Replicas: in.Replicas,
	}
	s.scaleChan <- *deployment
	return &commonTypes.Empty{}, nil
}

func (s *Server) UpdateManifest(ctx context.Context, in *commonTypes.Deployment) (*commonTypes.Empty, error) {
	deployment := &types.Deployment{
		ID:       in.ID,
		Manifest: in.Manifest,
	}
	s.reInitChan <- *deployment
	return &commonTypes.Empty{}, nil
}
