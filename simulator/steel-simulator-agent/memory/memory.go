package memory

import (
	"errors"
	"steel-lang/datastructure"
)

var (
	// ErrInvalidController represents an invalid or unknown controller
	ErrInvalidController error = errors.New("invalid controller")
	// ErrInvalidValues represents an invalid set of initialization values
	ErrInvalidValues error = errors.New("invalid initialization values")
)

// New creates a new memory, based on the passed memory controller and items
func New(controller string, items map[string]map[string][]string) (datastructure.ResourceController, error) {
	// I check the controller type and I return the correct implementation
	switch controller {
	case "basic":
		return NewBasicMemory(items)
	default:
		// If an invalid controller is passed, I raise an error
		return nil, ErrInvalidController
	}
}
