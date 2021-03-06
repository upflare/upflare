package upflare

import (
	"encoding/base64"
	"github.com/PuerkitoBio/purell"
	"net/url"
	"regexp"
)

var (
	FILENAME_PATTERN = regexp.MustCompile(`^[A-Za-z0-9\-_]*(\.[A-Za-z0-9\-_]+)+$`)
	HASH_PATTERN     = regexp.MustCompile(`^[0-9a-f]{40}$`)
	UPLOAD_PATTERN   = regexp.MustCompile(`^[A-Za-z0-9]{32}$`)
)

type Download struct {
	*App

	data string

	resize   *Resize
	filename string
}

func (d *Download) Resize(resize *Resize) *Download {
	d.resize = resize
	return d
}

func (d *Download) Filename(filename string) *Download {
	if FILENAME_PATTERN.MatchString(filename) {
		d.filename = filename
	}
	return d
}

func (d *Download) String() string {
	data := d.data
	if d.resize != nil {
		resize := d.resize.String()
		if resize != "" {
			data += "," + d.resize.String()
		}
	}
	if d.filename != "" {
		data += "," + d.filename
	}
	return d.Do("GET", "", data)
}

func (a *App) DownloadURL(u string) (*Download, error) {
	normalized, err := purell.NormalizeURLString(u, purell.FlagsSafe|purell.FlagRemoveDotSegments)
	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(normalized)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, ErrInvalidURL
	}

	if parsed.User != nil || parsed.Opaque != "" || parsed.Host == "" || parsed.Path == "" || parsed.Path[0] != '/' {
		return nil, ErrInvalidURL
	}

	parsed.Fragment = ""

	return &Download{
		App:  a,
		data: base64.URLEncoding.EncodeToString([]byte(parsed.String())),
	}, nil
}

func (a *App) DownloadHash(hash string) (*Download, error) {
	if !HASH_PATTERN.MatchString(hash) {
		return nil, ErrInvalidHash
	}
	return &Download{
		App:  a,
		data: hash,
	}, nil
}

func (a *App) DownloadUpload(id string) (*Download, error) {
	if !UPLOAD_PATTERN.MatchString(id) {
		return nil, ErrInvalidUpload
	}
	return &Download{
		App:  a,
		data: id,
	}, nil
}
