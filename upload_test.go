package upflare

import (
	"fmt"
	"testing"
)

var uploads = []struct {
	key    string
	secret string

	limit int64
	mimes []string

	width  int64
	height int64
	crop   bool

	timestamp int64
	nonce     string

	path string
	err  error
}{
	{
		key:    "key",
		secret: "secret",
		path:   "/key/eyJ0aW1lc3RhbXAiOjAsIm5vbmNlIjoiIn0=/-uL2Pwqn",
	},

	{
		key:    "key",
		secret: "secret",
		limit:  5 << 20,
		path:   "/key/eyJ0aW1lc3RhbXAiOjAsIm5vbmNlIjoiIiwibGltaXQiOjUyNDI4ODB9/xidLPwFG",
	},

	{
		key:    "key",
		secret: "secret",
		limit:  5 << 20,
		height: 60,
		path:   "/key/eyJ0aW1lc3RhbXAiOjAsIm5vbmNlIjoiIiwibGltaXQiOjUyNDI4ODAsInJlc2l6ZSI6Img2MCJ9/R7psRNmJ",
	},

	{
		key:    "key",
		secret: "secret",
		limit:  5 << 20,
		width:  60,
		crop:   true,
		mimes:  []string{"image/jpeg", "image/png"},
		path:   "/key/eyJ0aW1lc3RhbXAiOjAsIm5vbmNlIjoiIiwibGltaXQiOjUyNDI4ODAsImNvbnRlbnRfdHlwZXMiOlsiaW1hZ2UvanBlZyIsImltYWdlL3BuZyJdLCJyZXNpemUiOiJ3NjAtYyJ9/jx0Vzf8M",
	},
}

func TestUpload(t *testing.T) {
	for _, c := range uploads {
		a := New(c.key, c.secret)

		u, err := a.Upload()
		if err != c.err {
			t.Fatal(fmt.Sprintf("upflare: unexpected upload error: %v, expecting: %v", err, c.err))
		}
		if err != nil {
			continue
		}

		u.timestamp = c.timestamp
		u.nonce = c.nonce

		if c.width != 0 || c.height != 0 || c.crop {
			u.Resize(NewResize(c.width, c.height, c.crop))
		}

		if c.limit != 0 {
			u.Limit(c.limit)
		}

		if c.mimes != nil {
			u.ContentTypes(c.mimes...)
		}

		p := u.String()
		if p != c.path {
			t.Fatal(fmt.Sprintf("upflare: unexpected upload path: %s, expecting: %s", p, c.path))
		}
	}
}
