package api

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/luigi-riefolo/nlp/svc-storer/pb"
)

// StoreEntryHandler saves the received user entry.
func (s *StorerService) StoreEntryHandler(ctx context.Context, req *pb.StoreEntryRequest) (*empty.Empty, error) {
	log.WithField("id", req.Entry.Id).Debug("StoreEntryHandler")

	s.userEntries[req.Entry.Id] = req.Entry

	return &empty.Empty{}, nil
}
