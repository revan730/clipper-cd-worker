package api

import (
	"context"
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

func (s *Server) GetRevisions(ctx context.Context, in *commonTypes.RevisionsQuery) (*commonTypes.RevisionsArray, error) {
	revisions, err := s.databaseClient.FindRevisions(in.DeploymentID, in.Page, in.Limit)
	if err != nil {
		s.log.Error("Find revisions error", err)
		return &commonTypes.RevisionsArray{}, status.New(http.StatusInternalServerError, "").Err()
	}
	protoRevisions := &commonTypes.RevisionsArray{}
	for _, revision := range revisions {
		protoRevision := revisionToProto(revision)
		protoRevisions.Revisions = append(protoRevisions.Revisions, protoRevision)
	}
	return protoRevisions, nil
}
