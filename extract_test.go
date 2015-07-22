package upflare

import (
	"fmt"
	"testing"
)

var extracts = []struct {
	key    string
	secret string

	url string

	timestamp int64
	nonce     string

	path string
	err  error
}{
	{
		key:    "key",
		secret: "secret",
		url:    "http://www.google.com/",
		path:   "/extract/key/eyJ0aW1lc3RhbXAiOjAsIm5vbmNlIjoiIiwidXJsIjoiaHR0cDovL3d3dy5nb29nbGUuY29tLyJ9/2SZyJ_gk",
	},

	{
		key:    "key",
		secret: "secret",
		url:    "https://www.google.com/",
		path:   "/extract/key/eyJ0aW1lc3RhbXAiOjAsIm5vbmNlIjoiIiwidXJsIjoiaHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS8ifQ==/_8qejgJU",
	},

	{
		key:    "key",
		secret: "secret",
		url:    "https://www.google.com/robots.txt",
		path:   "/extract/key/eyJ0aW1lc3RhbXAiOjAsIm5vbmNlIjoiIiwidXJsIjoiaHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS9yb2JvdHMudHh0In0=/XYykQJm0",
	},

	{
		key:    "key",
		secret: "secret",
		url:    "https://www.google.com/ncr?a=1",
		path:   "/extract/key/eyJ0aW1lc3RhbXAiOjAsIm5vbmNlIjoiIiwidXJsIjoiaHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS9uY3I_YT0xIn0=/ztNX6sB8",
	},
}

func TestExtract(t *testing.T) {
	for _, c := range extracts {
		a := New(c.key, c.secret)

		u, err := a.Extract(c.url)
		if err != c.err {
			t.Fatal(fmt.Sprintf("upflare: unexpected extract error: %v, expecting: %v", err, c.err))
		}
		if err != nil {
			continue
		}

		u.timestamp = c.timestamp
		u.nonce = c.nonce

		p := u.String()
		if p != c.path {
			t.Fatal(fmt.Sprintf("upflare: unexpected extract path: %s, expecting: %s", p, c.path))
		}
	}
}
