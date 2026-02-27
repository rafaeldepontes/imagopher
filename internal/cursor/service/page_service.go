package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"os"
	"strconv"

	cursorMdl "github.com/rafaeldepontes/imagopher/internal/cursor/model"
)

type cursorService struct {
	secretKey []byte
	sigLen    int
}

func NewService() *cursorService {
	sigLenStr := os.Getenv("SIGNATURE_LENGTH")
	sigLen, _ := strconv.Atoi(sigLenStr)

	return &cursorService{
		secretKey: []byte(os.Getenv("SECRET_CURSOR_KEY")),
		sigLen:    sigLen,
	}
}

// Encode accepts a generic T type, a slice of any data, a size of records per page and
// the next page being a pointer to the next id in the database and it will return a hash
// with all the information needed in the next request for security.
func (s *cursorService) Encode(size int, nextCursor *uint64) (string, error) {
	rawData := cursorMdl.CursorBody{
		Size:       size,
		NextCursor: nextCursor,
	}

	sb, err := json.Marshal(rawData)
	if err != nil {
		return "", err
	}

	secretKey := os.Getenv("SECRET_CURSOR_KEY")
	var mac hash.Hash = hmac.New(sha256.New, []byte(secretKey))
	mac.Write(sb)
	signature := mac.Sum(nil)

	combined := append(sb, signature...)

	return base64.RawURLEncoding.EncodeToString(combined), nil
}

// Decode accepts a hashed source to decode, it will return the CursorPagination with the
// T type generic specified previously and an error if any.
func (s *cursorService) Decode(src string) (*cursorMdl.CursorBody, error) {
	combined, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return nil, err

	}
	if len(combined) < s.sigLen {
		return nil, errors.New("Invalid cursor length")
	}

	jsonBody := combined[:len(combined)-s.sigLen]
	signature := combined[len(combined)-s.sigLen:]

	var mac hash.Hash = hmac.New(sha256.New, s.secretKey)
	mac.Write(jsonBody)
	expected := mac.Sum(nil)

	if !hmac.Equal(signature, expected) {
		return nil, errors.New("Invalid cursor signature")
	}

	cursorModel := new(cursorMdl.CursorBody)
	_ = json.Unmarshal(jsonBody, cursorModel)

	return cursorModel, nil
}
