package model

import "errors"

var (
	QueryVariableNotFound = "%v query variable not found"
	IncorrectData         = "%v incorrect data"
	AlreadyExists         = "%v already exists"
)

var (
	ErrorRecordNotFound = errors.New("record not found")
	ErrorAlreadyExists  = errors.New("record already exists")
)
