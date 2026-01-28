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

	Resize    *ResizeOptions `json:"resize,omitempty"`
	Crop      *CropOptions   `json:"crop,omitempty"`
	Compress  *string        `json:"compress,omitempty"`
	Watermark *bool          `json:"watermark,omitempty"`
	Mirror    *bool          `json:"mirror,omitempty"`
	Rotate    *int           `json:"rotate,omitempty"`
	Format    *string        `json:"format,omitempty"`
	Filters   *FilterOptions `json:"filters,omitempty"`
}

type ResizeOptions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type CropOptions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type FilterOptions struct {
	Grayscale bool `json:"grayscale,omitempty"`
	Sepia     bool `json:"sepia,omitempty"`
}

type ImageResp struct {
	UUID uuid.UUID `json:"id"`
}
