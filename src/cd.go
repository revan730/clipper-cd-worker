package src

import (
	"github.com/golang/protobuf/proto"
	commonTypes "github.com/revan730/clipper-common/types"
	"go.uber.org/zap"
)

func (w *Worker) executeCDJob(CDJob commonTypes.CDJob) {
	w.logInfo("Got CD job message")
}

func (w *Worker) startConsuming() {
	defer w.jobsQueue.Close()
	blockMain := make(chan bool)

	cdMsgs, err := w.jobsQueue.MakeCDMsgChan()
	if err != nil {
		w.logFatal("Failed to create CD jobs channel", err)
	}

	go func() {
		for m := range cdMsgs {
			w.logger.Info("Received message: ", zap.ByteString("body", m))
			jobMsg := commonTypes.CDJob{}
			err := proto.Unmarshal(m, &jobMsg)
			if err != nil {
				w.logError("Failed to unmarshal job message", err)
				continue
			}
			go w.executeCDJob(jobMsg)
		}
	}()

	w.logInfo("Worker started")
	<-blockMain
}
