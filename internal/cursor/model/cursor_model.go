package model

// CursorBody is expected to have X size, the default should be 10.
type CursorBody struct {
	Size       int     `json:"size"`
	NextCursor *int64 `json:"next_cursor"`
}
