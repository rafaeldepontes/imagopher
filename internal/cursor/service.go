package cursor

import cursorMdl "github.com/rafaeldepontes/imagopher/internal/cursor/model"

type Page interface {
	Encode(size int, nextCursor *int64) (string, error)
	Decode(src string) (*cursorMdl.CursorBody, error)
}
