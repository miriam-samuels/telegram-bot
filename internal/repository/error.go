package types

import "fmt"

type CustomError struct {
	Code    int
	Message string
}

// Error implements error.
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
