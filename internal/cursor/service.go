package cursor

type Page interface {
	Encode(size int, nextCursor uint64) (string, error)
	Decode(src string) (any, error)
}
