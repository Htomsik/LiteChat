package dto

import (
	"time"
)

// MessageDTO Client version of Message
type MessageDTO struct {
	Type         string      `json:"type"`
	User         string      `json:"user"`
	Message      interface{} `json:"message"`
	DateTime     time.Time   `json:"dateTime"`
	clearPrivacy bool
}
