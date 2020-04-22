package errors

import (
	"errors"
	"runtime"
	"strings"
)

type (
	basicError struct {
		err error
		msg string
		pkg string
		fn  string
	}
)

func (e *basicError) Error() string {
	return e.msg
}

func (e *basicError) Unwrap() error {
	return e.err
}

// New returns an error that formats as the given text
func New(text string) error {
	p, f := packageAndFunc()
	return &basicError{err: errors.New(text), msg: text, pkg: p, fn: f}
}

// NewError wraps an error with additional metadata
func NewError(e error) error {
	p, f := packageAndFunc()
	return &basicError{err: e, msg: e.Error(), pkg: p, fn: f}
}

// see https://stackoverflow.com/questions/25262754/how-to-get-name-of-current-package-in-go
func packageAndFunc() (string, string) {
	pc, _, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	pkg := ""
	fn := parts[pl-1]
	if parts[pl-2][0] == '(' {
		fn = parts[pl-2] + "." + fn
		pkg = strings.Join(parts[0:pl-2], ".")
	} else {
		pkg = strings.Join(parts[0:pl-1], ".")
	}
	return pkg, fn
}
