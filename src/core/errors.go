package core

import (
	"errors"
)

var ErrNotPermittedOnline = errors.New(
	"Operation not permitted while core is online")
var ErrNotPermittedOffline = errors.New(
	"Operation not permitted while core is offline")
var ErrDuplicateServiceName = errors.New(
	"Duplicate service name")
