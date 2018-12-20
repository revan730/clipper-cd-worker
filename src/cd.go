package src

import (
	"github.com/golang/protobuf/proto"
	"github.com/revan730/clipper-cd-worker/types"
	commonTypes "github.com/revan730/clipper-common/types"
	"go.uber.org/zap"
)

// executeCDJob rolls new image onto k8s deployment
func (w *Worker) executeCDJob(CDJob commonTypes.CDJob) {
	w.logInfo("Got CD job message")
	// TODO: Get artifact gcr url using CI API
}

// initDeployment creates new deployment in k8s using manifest and
// provided image url
func (w *Worker) initDeployment(d types.Deployment) {
	w.logInfo("Initializing new deployment")
	// TODO: Get artifact gcr url using CI API
}

func (w *Worker) startConsuming() {
	defer w.jobsQueue.Close()
	blockMain := make(chan bool)

	cdMsgsQueue, err := w.jobsQueue.MakeCDMsgChan()
	if err != nil {
		w.logFatal("Failed to create CD jobs channel", err)
	}
	cdApiChan := w.apiServer.GetDepsChan()

	go func() {
		for {
			select {
			case m := <-cdMsgsQueue:
				w.logger.Info("Received message from queue: ", zap.ByteString("body", m))
				jobMsg := commonTypes.CDJob{}
				err := proto.Unmarshal(m, &jobMsg)
				if err != nil {
					w.logError("Failed to unmarshal job message", err)
					break
				}
				go w.executeCDJob(jobMsg)
			case m := <-cdApiChan:
				w.logInfo("New deployment: " + m.K8SName)
				go w.initDeployment(m)
			default:
				continue
			}
		}
	}()

	w.logInfo("Worker started")
	<-blockMain
}
