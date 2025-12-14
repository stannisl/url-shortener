package domain

import "errors"

var (
	ModelNotFoundErr = errors.New("model not found")
	NotUniqueErr     = errors.New("not unique")
)
