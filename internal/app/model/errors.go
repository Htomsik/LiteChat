package model

import "errors"

var (
	RecordNotFound           = errors.New("record not found")
	EmailOrPasswordIncorrect = errors.New("incorrect email or password")
	NotAuthenticated         = errors.New("not authenticated")
	NotActive                = errors.New("not active")
	ContextNotFound          = errors.New("user context not found")
)
