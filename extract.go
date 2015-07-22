package upflare

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/PuerkitoBio/purell"
	"io"
	"net/url"
	"time"
)

type Extract struct {
	*App

	url string

	timestamp int64
	nonce     string
}

func (u *Extract) String() string {
	data := &struct {
		Timestamp int64  `json:"timestamp"`
		Nonce     string `json:"nonce"`
		URL       string `json:"url"`
	}{
		Timestamp: u.timestamp,
		Nonce:     u.nonce,
		URL:       u.url,
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return u.Do("POST", "extract", base64.URLEncoding.EncodeToString(encoded))
}

func (a *App) Extract(u string) (*Extract, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, ErrInvalidURL
	}

	if parsed.User != nil || parsed.Opaque != "" || parsed.Host == "" || parsed.Path == "" || parsed.Path[0] != '/' || parsed.Fragment != "" {
		return nil, ErrInvalidURL
	}

	normalized, err := purell.NormalizeURLString(parsed.String(), purell.FlagsSafe)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	return &Extract{
		App:       a,
		url:       normalized,
		timestamp: time.Now().Unix(),
		nonce:     base64.URLEncoding.EncodeToString(nonce),
	}, nil
}
