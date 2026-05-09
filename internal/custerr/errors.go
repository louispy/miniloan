package custerr

import "errors"

var ErrDataNotFound = errors.New("data not found")
var ErrNoTransaction = errors.New("no open transaction")
