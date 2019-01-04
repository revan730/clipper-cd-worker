package db

import "github.com/revan730/clipper-cd-worker/types"

// DatabaseClient provides interface for data access layer operations
type DatabaseClient interface {
	Close()
	CreateSchema() error
	CreateDeployment(kd *types.Deployment) error
	CreateRevision(r *types.Revision) error
	SaveDeployment(kd *types.Deployment) error
}
