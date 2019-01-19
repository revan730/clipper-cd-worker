package src

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/revan730/clipper-cd-worker/types"
	commonTypes "github.com/revan730/clipper-common/types"
)

func renderManifestTemplate(manifest string, params types.ManifestValues) (string, error) {
	tpl := template.New("k8s manifest user's template")
	tpl, err := tpl.Parse(manifest)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, params)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (w *Worker) recordRevision(d types.Deployment, stdout string) error {
	revision := &types.Revision{
		DeploymentID: d.ID,
		ArtifactID:   d.ArtifactID,
		Date:         time.Now(),
		Stdout:       stdout,
		Replicas:     d.Replicas,
	}
	return w.databaseClient.CreateRevision(revision)
}

// updateDeploymentImage calls kubectl to change deployment image, using
// lock to syncronize update operations on deployment
func (w *Worker) updateDeploymentImage(dep types.Deployment, artifactID int64) {
	artifact, err := w.ciClient.GetBuildArtifactByID(dep.ArtifactID)
	if err != nil {
		w.log.Error("Failed to get build artifact", err)
		return
	}
	w.log.Info("Got artifact with url " + artifact.Name)
	lockRes := strconv.FormatInt(dep.ID, 10)
	err = w.distLock.Lock(lockRes)
	if err != nil {
		w.log.Error("Failed to acquire deployment lock", err)
		return
	}
	changedOk, stdout := w.kubectl.ChangeImage(dep.K8SName, artifact.Name)
	err = w.distLock.Unlock(lockRes)
	if err != nil {
		w.log.Error("Failed to release deployment lock", err)
	}
	dep.ArtifactID = artifact.ID
	err = w.recordRevision(dep, stdout)
	if err != nil {
		w.log.Error("Failed to record revision", err)
	}
	if changedOk == true {
		err = w.databaseClient.SaveDeployment(&dep)
		if err != nil {
			w.log.Error("Failed to save deployment to db", err)
			return
		}
	}
}

// scaleDeployment calls kubectl to scale deployment, using
// lock to syncronize update operations on deployment
func (w *Worker) scaleDeployment(dep types.Deployment, replicas int64) {
	lockRes := strconv.FormatInt(dep.ID, 10)
	err := w.distLock.Lock(lockRes)
	if err != nil {
		w.log.Error("Failed to acquire deployment lock", err)
		return
	}
	scaledOk, stdout := w.kubectl.ScaleDeployment(dep.K8SName, replicas)
	err = w.distLock.Unlock(lockRes)
	if err != nil {
		w.log.Error("Failed to release deployment lock", err)
	}
	dep.Replicas = replicas
	err = w.recordRevision(dep, stdout)
	if err != nil {
		w.log.Error("Failed to record revision", err)
	}
	if scaledOk == true {
		err = w.databaseClient.SaveDeployment(&dep)
		if err != nil {
			w.log.Error("Failed to save deployment to db", err)
			return
		}
	}
}

// executeCDJob rolls new image onto k8s deployment
func (w *Worker) executeCDJob(CDJob commonTypes.CDJob) {
	w.log.Info("Got CD job message")
	deployments, err := w.databaseClient.FindDeploymentsByRepo(CDJob.RepoID)
	if err != nil {
		w.log.Error("Failed to execute CD job", err)
		return
	}
	for _, dep := range deployments {
		go w.updateDeploymentImage(dep, CDJob.ArtifactID)
	}
}

func (w *Worker) updateImageFromProto(d types.Deployment) {
	// Find deployment in database and call image update
	deployment, err := w.databaseClient.FindDeployment(d.ID)
	if err != nil {
		w.log.Error("Failed to find deployment", err)
		return
	}
	go w.updateDeploymentImage(*deployment, d.ArtifactID)
}

func (w *Worker) scaleFromProto(d types.Deployment) {
	// Find deployment in database and call deployment scale
	deployment, err := w.databaseClient.FindDeployment(d.ID)
	if err != nil {
		w.log.Error("Failed to find deployment", err)
		return
	}
	go w.scaleDeployment(*deployment, d.Replicas)
}

// initDeployment creates new deployment in k8s using manifest and
// provided image url
func (w *Worker) initDeployment(d types.Deployment) {
	w.log.Info("Initializing new deployment")
	artifact, err := w.ciClient.GetBuildArtifactByID(d.ArtifactID)
	if err != nil {
		w.log.Error("Failed to get build artifact", err)
		return
	}
	manifestVals := types.ManifestValues{
		Name:     d.K8SName,
		Image:    artifact.Name,
		Replicas: d.Replicas,
	}
	manifest, err := renderManifestTemplate(d.Manifest, manifestVals)
	if err != nil {
		w.log.Error("Failed to render manifest template", err)
		return
	}
	fmt.Println("Manifest:\n" + manifest)
	ok, stdout := w.kubectl.CreateDeployment(manifest)
	if ok != true {
		fmt.Println("fucked up")
	}
	fmt.Println("stdout: " + stdout)
	err = w.recordRevision(d, stdout)
	if err != nil {
		w.log.Error("Failed to write revision to db", err)
	}
	if ok == true {
		d.IsInitialized = true
		err = w.databaseClient.SaveDeployment(&d)
		if err != nil {
			w.log.Error("Failed to update deployment db record", err)
		}
	}
}

func (w *Worker) startConsuming() {
	defer w.jobsQueue.Close()
	blockMain := make(chan bool)

	cdMsgsQueue, err := w.jobsQueue.MakeCDMsgChan()
	if err != nil {
		w.log.Fatal("Failed to create CD jobs channel", err)
	}
	cdAPIChan := w.apiServer.GetDepsChan()
	changeImageChan := w.apiServer.GetImageChangeChan()
	scaleChan := w.apiServer.GetScaleChan()

	go func() {
		for {
			select {
			case m := <-cdMsgsQueue:
				body := string(m)
				w.log.Info("Received message from queue: " + body)
				jobMsg := commonTypes.CDJob{}
				err := proto.Unmarshal(m, &jobMsg)
				if err != nil {
					w.log.Error("Failed to unmarshal job message", err)
					break
				}
				go w.executeCDJob(jobMsg)
			case m := <-cdAPIChan:
				w.log.Info("New deployment: " + m.K8SName)
				go w.initDeployment(m)
			case m := <-changeImageChan:
				go w.updateImageFromProto(m)
			case m := <-scaleChan:
				go w.scaleFromProto(m)
			}
		}
	}()

	w.log.Info("Worker started")
	<-blockMain
}
