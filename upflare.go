package upflare

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
)

type App struct {
	key, secret string
}

func New(key, secret string) *App {
	return &App{
		key:    key,
		secret: secret,
	}
}

func (a *App) Sign(base []byte) string {
	mac := hmac.New(sha1.New, []byte(a.secret))
	mac.Write(base)
	signature := mac.Sum(nil)
	return base64.URLEncoding.EncodeToString(signature[len(signature)-6:])
}

func (a *App) Do(method, data string) string {
	signature := a.Sign([]byte(fmt.Sprintf("%s /%s/%s", strings.ToUpper(method), a.key, data)))
	return fmt.Sprintf("/%s/%s/%s", a.key, data, signature)
}

type Resize struct {
	width  int64
	height int64
	crop   bool
}

func (r *Resize) String() string {
	resize := ""
	switch {
	case r.width > 0 && r.height > 0 && r.width == r.height:
		resize = fmt.Sprintf("s%d", r.width)
	case r.width > 0 && r.height > 0:
		resize = fmt.Sprintf("w%d-h%d", r.width, r.height)
	case r.width > 0:
		resize = fmt.Sprintf("w%d", r.width)
	case r.height > 0:
		resize = fmt.Sprintf("h%d", r.height)
	}
	if resize != "" && r.crop {
		resize += "-c"
	}
	return resize
}

func NewResize(width, height int64, crop bool) *Resize {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	if width == 0 && height == 0 && !crop {
		return nil
	}
	return &Resize{
		width:  width,
		height: height,
		crop:   crop,
	}
}
