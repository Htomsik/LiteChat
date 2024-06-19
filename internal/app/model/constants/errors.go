package constants

import "errors"

var (
	QueryVariableNotFound = "%v query variable not found"
)

var (
	ErrorRecordNotFound = errors.New("record not found")
	ErrorAlreadyExists  = errors.New("record already exists")
)
