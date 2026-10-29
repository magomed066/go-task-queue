package job

import (
	"encoding/json"
	"time"
)

type Type string

const (
	TypeSendEmail   Type = "send_email"
	TypeResizeImage Type = "resize_image"
)

type Task struct {
	ID        string
	Type      Type
	Payload   json.RawMessage
	CreatedAt time.Time
	Timeout   time.Duration
}
