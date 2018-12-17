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
		(*types.Deployment)(nil)} {
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
