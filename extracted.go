package upflare

import (
	"encoding/json"
	"math"
	"time"
)

type Extracted struct {
	Values map[string]interface{}
}

func (a *App) Extracted(data []byte) (*Extracted, error) {
	meta := &struct {
		Key       string `json:"key"`
		Timestamp int64  `json:"timestamp"`
		Nonce     string `json:"nonce"`
		Signature string `json:"signature"`
	}{}

	if err := json.Unmarshal(data, meta); err != nil {
		return nil, err
	}

	if meta.Key != a.key {
		return nil, ErrInvalidSignature
	}

	if math.Abs(time.Since(time.Unix(meta.Timestamp, 0)).Minutes()) > 15 {
		return nil, ErrInvalidSignature
	}

	var unknown interface{}
	if err := json.Unmarshal(data, &unknown); err != nil {
		return nil, err
	}

	values, ok := unknown.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidSignature
	}

	delete(values, "signature")
	values["timestamp"] = meta.Timestamp

	base, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	if a.Sign(base) != meta.Signature {
		return nil, ErrInvalidSignature
	}

	delete(values, "timestamp")
	delete(values, "key")
	delete(values, "nonce")

	return &Extracted{
		Values: values,
	}, nil
}
