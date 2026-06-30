package validator

import (
	"errors"
	"strings"
)

func ValidateTarget(target string) error {

	target = strings.TrimSpace(target)

	if target == "" {
		return errors.New("target cannot be empty")
	}

	return nil
}
