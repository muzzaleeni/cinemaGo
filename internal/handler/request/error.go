package request

import "fmt"

type BlankFieldError struct {
	field string
}

func (e BlankFieldError) Error() string {
	return fmt.Sprintf("%q can not be blank", e.field)
}
