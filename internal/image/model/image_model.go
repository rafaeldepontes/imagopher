package model

import "github.com/google/uuid"

type ImageEntity struct {
	ID   uint64    `json:"id"`
	Path string    `json:"path"`
	UUID uuid.UUID `json:"uuid"`
}

// TransformReq has an UUID as their "id", this is a protection
// to not expose the REAL id from the database...
type TransformReq struct {
	UUID uuid.UUID `json:"id"`

	// TODO: Add the transformation... for now I only have the UUID
}

type ImageResp struct {
	UUID uuid.UUID `json:"id"`
}
