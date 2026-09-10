package ipmi

import (
	"errors"
	"fmt"
)

type exitCodes map[int]string

func (codes exitCodes) getErrorMessage(code int) error {
	if codeInfo, ok := codes[code]; ok {
		return errors.New(codeInfo)
	}
	return fmt.Errorf("exit with status: %d", code)
}
