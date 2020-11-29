package ecs

import (
	"errors"
	"fmt"
)

var (
	errMissingComponent = fmt.Errorf("missing component")
	errMissingEntity    = fmt.Errorf("missing entity")
)

func IsMissingComponentError(err error) bool {
	return errors.Is(err, errMissingComponent)
}

func IsMissingEntityError(err error) bool {
	return errors.Is(err, errMissingEntity)
}
