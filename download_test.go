package upflare

import (
	"fmt"
	"testing"
)

var downloads = []struct {
	key    string
	secret string

	url    string
	hash   string
	upload string

	width  int64
	height int64
	crop   bool

	filename string

	path string
	err  error
}{
	{
		key:    "key",
		secret: "secret",
		url:    "http://www.google.com/favicon.ico",
		path:   "/key/aHR0cDovL3d3dy5nb29nbGUuY29tL2Zhdmljb24uaWNv/IDia8aMJ",
	},

	{
		key:    "key",
		secret: "secret",
		url:    "http://www.google.com/favicon.ico",
		width:  60,
		height: 60,
		path:   "/key/aHR0cDovL3d3dy5nb29nbGUuY29tL2Zhdmljb24uaWNv,s60/1LUdHI96",
	},

	{
		key:    "key",
		secret: "secret",
		url:    "http://www.google.com/favicon.ico",
		width:  60,
		height: 60,
		crop:   true,
		path:   "/key/aHR0cDovL3d3dy5nb29nbGUuY29tL2Zhdmljb24uaWNv,s60-c/NoZBbIeE",
	},

	{
		key:      "key",
		secret:   "secret",
		url:      "http://www.google.com/favicon.ico",
		width:    60,
		height:   60,
		crop:     true,
		filename: "test.ico",
		path:     "/key/aHR0cDovL3d3dy5nb29nbGUuY29tL2Zhdmljb24uaWNv,s60-c,test.ico/1jGcZmFI",
	},

	{
		key:      "key",
		secret:   "secret",
		url:      "http://www.google.com/favicon.ico",
		crop:     true,
		filename: "test.ico",
		path:     "/key/aHR0cDovL3d3dy5nb29nbGUuY29tL2Zhdmljb24uaWNv,test.ico/hs205LRj",
	},

	{
		key:    "key",
		secret: "secret",
		url:    "ftp://www.google.com/favicon.ico",
		err:    ErrInvalidURL,
	},

	{
		key:    "key",
		secret: "secret",
		url:    "http://www.google.com/favicon.ico#test",
		err:    ErrInvalidURL,
	},

	{
		key:    "key",
		secret: "secret",
		hash:   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		path:   "/key/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/OYzWZrtL",
	},

	{
		key:    "key",
		secret: "secret",
		hash:   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaA",
		err:    ErrInvalidHash,
	},

	{
		key:    "key",
		secret: "secret",
		upload: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		path:   "/key/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/iRWIM9L2",
	},

	{
		key:    "key",
		secret: "secret",
		upload: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-",
		err:    ErrInvalidUpload,
	},
}

func TestDownload(t *testing.T) {
	for _, c := range downloads {
		a := New(c.key, c.secret)

		var d *Download
		var err error

		switch {
		case c.url != "":
			d, err = a.DownloadURL(c.url)
		case c.hash != "":
			d, err = a.DownloadHash(c.hash)
		case c.upload != "":
			d, err = a.DownloadUpload(c.upload)
		}

		if err != c.err {
			t.Fatal(fmt.Sprintf("upflare: unexpected download error: %v, expecting: %v", err, c.err))
		}

		if err != nil {
			continue
		}

		if c.width != 0 || c.height != 0 || c.crop {
			d.Resize(c.width, c.height, c.crop)
		}

		if c.filename != "" {
			d.Filename(c.filename)
		}

		p := d.String()
		if p != c.path {
			t.Fatal(fmt.Sprintf("upflare: unexpected download path: %s, expecting: %s", p, c.path))
		}

	}
}
