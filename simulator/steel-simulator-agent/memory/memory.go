package memory

import (
	"errors"
	"steel/memory"
)

// New creates a new memory, based on the passed memory controller and items
func New(controller string, items map[string]map[string]string) (memory.ResourceController, error) {
	// I check the controller type and I return the correct implementation
	switch controller {
	case "basic":
		return NewBasicMemory(items)
	default:
		// If an invalid controller is passed, I raise an error
		return nil, errors.New("invalid controller")
	}
}
