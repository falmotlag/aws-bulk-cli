package errors

import (
	"fmt"

	goerrors "github.com/go-errors/errors"
)

func WithStackTrace(err error) error {
	if err == nil {
		return nil
	}

	return goerrors.Wrap(err, 1)
}

func WithStackTraceAndPrefix(err error, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return goerrors.WrapPrefix(err, fmt.Sprintf(message, args...), 1)
}
