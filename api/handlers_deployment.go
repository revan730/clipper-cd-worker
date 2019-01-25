package api

import (
	"context"
	"fmt"
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

func (s *Server) DeleteDeployment(ctx context.Context, in *commonTypes.Deployment) (*commonTypes.Empty, error) {
	deployment, err := s.databaseClient.FindDeployment(in.ID)
	if err != nil {
		s.log.Error("Delete deployment error: couldn't find deployment", err)
		return &commonTypes.Empty{}, status.New(http.StatusInternalServerError, "").Err()
	}
	if deployment == nil {
		s.log.Info(fmt.Sprintf("Delete deployment - deployment not found with id %d", in.ID))
		return &commonTypes.Empty{}, status.New(http.StatusBadRequest, "").Err()
	}
	err = s.databaseClient.DeleteDeployment(deployment)
	if err != nil {
		s.log.Error("Delete deployment error - couldn't delete deployment", err)
		return &commonTypes.Empty{}, status.New(http.StatusInternalServerError, "").Err()
	}
	s.deleteChan <- *deployment
	return &commonTypes.Empty{}, nil
}
