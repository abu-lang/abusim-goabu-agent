package memory

import (
	"errors"
	"steel-lang/datastructure"
)

var (
	ErrInvalidController error = errors.New("invalid controller")
	ErrInvalidValues     error = errors.New("invalid initialization values")
)

func New(controller string, items map[string]map[string][]string) (datastructure.ResourceController, error) {
	switch controller {
	case "basic":
		return NewBasicMemory(items)
	default:
		return nil, ErrInvalidController
	}
}
