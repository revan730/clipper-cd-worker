package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/revan730/clipper-cd-worker/types"
)

// PostgresClient provides data access layer to objects in Postgres.
// implements DatabaseClient interface
type PostgresClient struct {
	pg *pg.DB
}

// NewPGClient creates new copy of PostgresClient
func NewPGClient(config types.PGClientConfig) *PostgresClient {
	DBClient := &PostgresClient{}
	pgdb := pg.Connect(&pg.Options{
		User:         config.DBUser,
		Addr:         config.DBAddr,
		Password:     config.DBPassword,
		Database:     config.DB,
		MinIdleConns: 2,
	})
	DBClient.pg = pgdb
	return DBClient
}

// Close gracefully closes db connection
func (d *PostgresClient) Close() {
	d.pg.Close()
}

// CreateSchema creates database tables if they not exist
func (d *PostgresClient) CreateSchema() error {
	for _, model := range []interface{}{
		(*types.Deployment)(nil),
		(*types.Revision)(nil)} {
		err := d.pg.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateDeployment creates k8s deployment record from provided struct
func (d *PostgresClient) CreateDeployment(kd *types.Deployment) error {
	return d.pg.Insert(kd)
}

// FindDeploymentsByRepo returns all deployments for provided repo id
func (d *PostgresClient) FindDeploymentsByRepo(repoID int64) ([]types.Deployment, error) {
	var deployments []types.Deployment

	err := d.pg.Model(&deployments).
		Where("repo_id = ?", repoID).
		Select()

	return deployments, err
}

// CreateRevision creates deployment revision record from provided struct
func (d *PostgresClient) CreateRevision(r *types.Revision) error {
	return d.pg.Insert(r)
}

// SaveDeployment updates provided deployment in db
func (d *PostgresClient) SaveDeployment(kd *types.Deployment) error {
	return d.pg.Update(kd)
}
