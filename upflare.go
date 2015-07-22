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

func (a *App) Do(method, action, data string) string {
	if action == "" {
		signature := a.Sign([]byte(fmt.Sprintf("%s /%s/%s", strings.ToUpper(method), a.key, data)))
		return fmt.Sprintf("/%s/%s/%s", a.key, data, signature)
	} else {
		signature := a.Sign([]byte(fmt.Sprintf("%s /%s/%s/%s", strings.ToUpper(method), action, a.key, data)))
		return fmt.Sprintf("/%s/%s/%s/%s", action, a.key, data, signature)
	}
}
