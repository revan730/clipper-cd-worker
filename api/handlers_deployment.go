package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/revan730/clipper-cd-worker/types"
	commonTypes "github.com/revan730/clipper-common/types"
	"google.golang.org/grpc/status"
)

func deploymentToProto(deployment *types.Deployment) *commonTypes.Deployment {
	return &commonTypes.Deployment{
		ID:         deployment.ID,
		Branch:     deployment.Branch,
		RepoID:     deployment.RepoID,
		ArtifactID: deployment.ArtifactID,
		K8SName:    deployment.K8SName,
		Manifest:   deployment.Manifest,
		Replicas:   deployment.Replicas,
	}
}

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

func (s *Server) GetAllDeployments(ctx context.Context, in *commonTypes.DeploymentsQuery) (*commonTypes.DeploymentsArray, error) {
	deployments, err := s.databaseClient.FindAllDeployments(in.Page, in.Limit)
	if err != nil {
		s.log.Error("Find all deployments error", err)
		return &commonTypes.DeploymentsArray{}, status.New(http.StatusInternalServerError, "").Err()
	}
	count, err := s.databaseClient.FindDeploymentCount()
	if err != nil {
		s.log.Error("Find deployments count error", err)
		return &commonTypes.DeploymentsArray{}, status.New(http.StatusInternalServerError, "").Err()
	}
	protoDeps := &commonTypes.DeploymentsArray{}
	for _, dep := range deployments {
		protoDep := deploymentToProto(dep)
		protoDeps.Deployments = append(protoDeps.Deployments, protoDep)
	}
	protoDeps.Total = count
	return protoDeps, nil
}

func (s *Server) GetDeployment(ctx context.Context, in *commonTypes.Deployment) (*commonTypes.Deployment, error) {
	deployment, err := s.databaseClient.FindDeployment(in.ID)
	if err != nil {
		s.log.Error("Get deployment error", err)
		return &commonTypes.Deployment{}, status.New(http.StatusInternalServerError, "").Err()
	}
	if deployment == nil {
		s.log.Info(fmt.Sprintf("Get deployment - deployment not found with id %d", in.ID))
		return &commonTypes.Deployment{}, status.New(http.StatusNotFound, "").Err()
	}
	return deploymentToProto(deployment), nil
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
