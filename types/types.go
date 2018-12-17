package types

type Deployment struct {
	ID            int64  `json:"deploymentID"`
	Branch        string `json:"branch"`
	ArtifactID    int64  `json:"artifactID"`
	K8SName       string `json:"k8sName"`
	Manifest      string `json:"manifest"`
	IsInitialized bool   `json:"isInitialized" sql:"default:false"`
}

type PGClientConfig struct {
	DBAddr     string
	DB         string
	DBUser     string
	DBPassword string
}
