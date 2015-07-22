package upflare

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"
)

type Upload struct {
	*App

	limit        int64
	contentTypes []string
	resize       *Resize

	timestamp int64
	nonce     string
}

func (u *Upload) ContentTypes(contentTypes ...string) *Upload {
	u.contentTypes = contentTypes
	return u
}

func (u *Upload) Limit(limit int64) *Upload {
	u.limit = limit
	return u
}

func (u *Upload) Resize(resize *Resize) *Upload {
	u.resize = resize
	return u
}

func (u *Upload) String() string {
	resize := ""
	if u.resize != nil {
		resize = u.resize.String()
	}

	data := &struct {
		Timestamp    int64    `json:"timestamp"`
		Nonce        string   `json:"nonce"`
		Limit        int64    `json:"limit,omitempty"`
		ContentTypes []string `json:"content_types,omitempty"`
		Resize       string   `json:"resize,omitempty"`
	}{
		Timestamp:    u.timestamp,
		Nonce:        u.nonce,
		Limit:        u.limit,
		ContentTypes: u.contentTypes,
		Resize:       resize,
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return u.Do("POST", "upload", base64.URLEncoding.EncodeToString(encoded))
}

func (a *App) Upload() (*Upload, error) {

	nonce := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	return &Upload{
		App:       a,
		timestamp: time.Now().Unix(),
		nonce:     base64.URLEncoding.EncodeToString(nonce),
	}, nil
}
