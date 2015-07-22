package upflare

import (
	"encoding/json"
	"errors"
	"math"
	"time"
)

var (
	ErrInvalidSignature = errors.New("upflare: invalid signature")
)

type Uploaded struct {
	ID string

	Filename    string
	Hash        string
	Size        int64
	ContentType string

	Width  int64
	Height int64
}

func (a *App) Uploaded(data []byte) (*Uploaded, error) {
	value := &struct {
		Key         string `json:"key"`
		ID          string `json:"id"`
		Filename    string `json:"filename"`
		Hash        string `json:"hash"`
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
		Width       int64  `json:"width,omitempty"`
		Height      int64  `json:"height,omitempty"`
		Timestamp   int64  `json:"timestamp"`
		Nonce       string `json:"nonce"`
		Signature   string `json:"signature"`
	}{}

	if err := json.Unmarshal(data, value); err != nil {
		return nil, err
	}

	if value.Key != a.key {
		return nil, ErrInvalidSignature
	}

	values := map[string]interface{}{
		"key":          value.Key,
		"id":           value.ID,
		"filename":     value.Filename,
		"hash":         value.Hash,
		"size":         value.Size,
		"content_type": value.ContentType,
		"timestamp":    value.Timestamp,
		"nonce":        value.Nonce,
	}
	if value.Width > 0 {
		values["width"] = value.Width
	}
	if value.Height > 0 {
		values["height"] = value.Height
	}

	base, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	if a.Sign(base) != value.Signature {
		return nil, ErrInvalidSignature
	}

	if math.Abs(time.Since(time.Unix(value.Timestamp, 0)).Minutes()) > 15 {
		return nil, ErrInvalidSignature
	}

	return &Uploaded{
		ID:          value.ID,
		Filename:    value.Filename,
		Hash:        value.Hash,
		Size:        value.Size,
		ContentType: value.ContentType,
		Width:       value.Width,
		Height:      value.Height,
	}, nil
}
