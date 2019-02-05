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
	FindDeploymentCount() (int64, error)
	FindDeploymentsByRepo(repoID int64) ([]types.Deployment, error)
	CreateRevision(r *types.Revision) error
	FindRevision(revisionID int64) (*types.Revision, error)
	FindRevisions(deploymentID, page, limit int64) ([]*types.Revision, error)
	FindRevisionsCount(deploymentID int64) (int64, error)
	SaveDeployment(kd *types.Deployment) error
}
