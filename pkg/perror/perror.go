package perror

import (
	"errors"
	"fmt"
)

func NewErrorString(f string, a... interface{}) error {
	return errors.New(fmt.Sprintf(f, a...))
}
