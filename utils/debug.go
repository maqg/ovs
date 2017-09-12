package utils

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

// Assertf for assert
func Assertf(expression bool, f string, args ...interface{}) {
	Assert(expression, fmt.Sprintf(f, args...))
}

// Assert Management
func Assert(expression bool, msg string) {
	if !expression {
		panic(errors.New(msg))
	}
}

// PanicIfError for panic management
func PanicIfError(ok bool, err error) {
	if !ok {
		panic(err)
	}
}

// PanicOnError for error management
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// LogError for error logging
func LogError(args ...interface{}) {
	for _, arg := range args {
		if e, ok := arg.(error); ok {
			err := errors.Wrap(e, "UNHANDLED ERROR, PLEASE REPORT A BUG TO US")
			log.Warn(fmt.Sprintf("%+v\n", err))
		}
	}
}
