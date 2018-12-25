package src

import (
	"bytes"
	"fmt"
	"text/template"

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

// executeCDJob rolls new image onto k8s deployment
func (w *Worker) executeCDJob(CDJob commonTypes.CDJob) {
	w.log.LogInfo("Got CD job message")
	// TODO: Get artifact gcr url using CI API
	// TODO: Call kubectl to change image
	// TODO: Record result to revisions
}

// initDeployment creates new deployment in k8s using manifest and
// provided image url
func (w *Worker) initDeployment(d types.Deployment) {
	w.log.LogInfo("Initializing new deployment")
	artifact, err := w.ciClient.GetBuildArtifactByID(d.ArtifactID)
	if err != nil {
		w.log.LogError("Failed to get build artifact", err)
	}
	manifestVals := types.ManifestValues{
		Name:     d.K8SName,
		Image:    artifact.Name,
		Replicas: d.Replicas,
	}
	manifest, err := renderManifestTemplate(d.Manifest, manifestVals)
	if err != nil {
		w.log.LogError("Failed to render manifest template", err)
	}
	fmt.Println("Manifest:\n" + manifest)
	// TODO: Call kubectl to create deployment
	// TODO: Record result to revisions
	// TODO: Change deployment isInitialized flag
}

func (w *Worker) startConsuming() {
	defer w.jobsQueue.Close()
	blockMain := make(chan bool)

	cdMsgsQueue, err := w.jobsQueue.MakeCDMsgChan()
	if err != nil {
		w.log.LogFatal("Failed to create CD jobs channel", err)
	}
	cdAPIChan := w.apiServer.GetDepsChan()

	go func() {
		for {
			select {
			case m := <-cdMsgsQueue:
				body := string(m)
				w.log.LogInfo("Received message from queue: " + body)
				jobMsg := commonTypes.CDJob{}
				err := proto.Unmarshal(m, &jobMsg)
				if err != nil {
					w.log.LogError("Failed to unmarshal job message", err)
					break
				}
				go w.executeCDJob(jobMsg)
			case m := <-cdAPIChan:
				w.log.LogInfo("New deployment: " + m.K8SName)
				go w.initDeployment(m)
			}
		}
	}()

	w.log.LogInfo("Worker started")
	<-blockMain
}
