package db

import "github.com/revan730/clipper-cd-worker/types"

// DatabaseClient provides interface for data access layer operations
type DatabaseClient interface {
	Close()
	CreateSchema() error
	CreateDeployment(kd *types.Deployment) error
	DeleteDeployment(kd *types.Deployment) error
	FindAllDeployments(page, limit int64) ([]*types.Deployment, error)
	FindDeployment(deploymentID int64) (*types.Deployment, error)
	FindDeploymentsByRepo(repoID int64) ([]types.Deployment, error)
	CreateRevision(r *types.Revision) error
	FindRevisions(deploymentID, page, limit int64) ([]*types.Revision, error)
	SaveDeployment(kd *types.Deployment) error
}
