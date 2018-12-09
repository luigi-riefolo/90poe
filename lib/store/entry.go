package store

import (
	storerpb "github.com/luigi-riefolo/90poe/svc-storer/pb"
)

// Entry represents a user record in a storer.
type Entry struct {
	ID string

	Name         string
	Email        string
	MobileNumber string

	CreateAt  int
	UpdatedAt int
}

// ToProto converts a storer entry to its protobuf message representation.
func (e *Entry) ToProto() storerpb.Entry {

	entry := storerpb.Entry{
		Id:           e.ID,
		Name:         e.Name,
		Email:        e.Email,
		MobileNumber: e.MobileNumber,
	}

	return entry
}

// FromProto converts an Entry protobuf message to its storer representation.
func FromProto(e *storerpb.Entry) *Entry {

	entry := &Entry{
		ID:           e.Id,
		Name:         e.Name,
		Email:        e.Email,
		MobileNumber: e.MobileNumber,
	}

	return entry
}
