package api

import (
	"context"
	"fmt"
	"net/http"

	ptypes "github.com/golang/protobuf/ptypes"
	"github.com/revan730/clipper-cd-worker/types"
	commonTypes "github.com/revan730/clipper-common/types"
	"google.golang.org/grpc/status"
)

func revisionToProto(r *types.Revision) *commonTypes.Revision {
	timestamp, _ := ptypes.TimestampProto(r.Date)
	return &commonTypes.Revision{
		ID:           r.ID,
		ArtifactID:   r.ArtifactID,
		DeploymentID: r.DeploymentID,
		Date:         timestamp,
		Stdout:       r.Stdout,
		Replicas:     r.Replicas,
	}
}

func (s *Server) GetRevision(ctx context.Context, in *commonTypes.Revision) (*commonTypes.Revision, error) {
	revision, err := s.databaseClient.FindRevision(in.ID)
	if err != nil {
		s.log.Error("Get revision error", err)
		return &commonTypes.Revision{}, status.New(http.StatusInternalServerError, "").Err()
	}
	if revision == nil {
		s.log.Info(fmt.Sprintf("Get revision - revision not found with id %d", in.ID))
		return &commonTypes.Revision{}, status.New(http.StatusNotFound, "").Err()
	}
	return revisionToProto(revision), nil
}

func (s *Server) GetRevisions(ctx context.Context, in *commonTypes.RevisionsQuery) (*commonTypes.RevisionsArray, error) {
	revisions, err := s.databaseClient.FindRevisions(in.DeploymentID, in.Page, in.Limit)
	if err != nil {
		s.log.Error("Find revisions error", err)
		return &commonTypes.RevisionsArray{}, status.New(http.StatusInternalServerError, "").Err()
	}
	count, err := s.databaseClient.FindRevisionsCount(in.DeploymentID)
	if err != nil {
		s.log.Error("Find revisions count error", err)
		return &commonTypes.RevisionsArray{}, status.New(http.StatusInternalServerError, "").Err()
	}
	protoRevisions := &commonTypes.RevisionsArray{}
	for _, revision := range revisions {
		protoRevision := revisionToProto(revision)
		protoRevisions.Revisions = append(protoRevisions.Revisions, protoRevision)
	}
	protoRevisions.Total = count
	return protoRevisions, nil
}
