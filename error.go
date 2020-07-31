package lazy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var (
	// ErrNoConfiguration ...
	ErrNoConfiguration = errors.New("can't find lazy configuration")
	// ErrHasAssociations ...
	ErrHasAssociations = errors.New("can't operate because of associations")
	// ErrUnknown ...
	ErrUnknown = errors.New("unknown")
)

// Errors contains all happened errors
type Errors []error

// GetErrors gets all happened errors
func (errs Errors) GetErrors() []error {
	return errs
}

// Add adds an error
func (errs Errors) Add(newErrors ...error) Errors {
	for _, err := range newErrors {
		if err == nil {
			continue
		}

		if errors, ok := err.(Errors); ok {
			errs = errs.Add(errors...)
		} else {
			ok = true
			for _, e := range errs {
				if err == e {
					ok = false
				}
			}
			if ok {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// Error format happened errors
func (errs Errors) Error() string {
	var errors = []string{}
	for _, e := range errs {
		errors = append(errors, e.Error())
	}
	return strings.Join(errors, "; ")
}

// WarpErr ...
func WarpErr(desc string, err error) error {
	return fmt.Errorf("%v --> %w", desc, err)
}

// ErrorCode ...
func ErrorCode(err error) (code uint, msg string) {
	var mySQLError *mysql.MySQLError
	switch {
	case errors.As(err, &mySQLError):
		if val, ok := err.(*mysql.MySQLError); ok {
			return uint(val.Number), val.Message
		}
	}
	return 0, "unknown error"
}
