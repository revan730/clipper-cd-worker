package types

import "time"

type Deployment struct {
	ID            int64  `json:"deploymentID"`
	RepoID        int64  `json:"repoID"`
	Branch        string `json:"branch"`
	ArtifactID    int64  `json:"artifactID"`
	K8SName       string `sql:",unique" json:"k8sName"`
	Manifest      string `json:"manifest"`
	Replicas      int64  `json:"replicas"`
	IsInitialized bool   `json:"isInitialized" sql:"default:false"`
}

type Revision struct {
	ID           int64     `json:"revisionID"`
	DeploymentID int64     `json:"deploymentID"`
	ArtifactID   int64     `json:"artifactID"`
	Date         time.Time `json:"date"`
	Stdout       string    `json:"stdout"`
	Replicas     int64     `json:"replicas"`
}

type PGClientConfig struct {
	DBAddr     string
	DB         string
	DBUser     string
	DBPassword string
}
