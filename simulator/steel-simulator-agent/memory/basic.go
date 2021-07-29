package memory

import (
	"fmt"
	"steel-lang/datastructure"
	"strconv"
	"time"
)

// New creates a new memory, based on the basic resources
func NewBasicMemory(items map[string]map[string]string) (datastructure.ResourceController, error) {
	// I create an empty basic memory...
	mem := datastructure.MakeResources()
	// ... and I range over the items to initialize it, with the provided initialization value or a default
	for vartype, values := range items {
		for name, initvalue := range values {
			switch vartype {
			case "bool":
				if val, err := getBasicMemoryBool(initvalue); err != nil {
					return nil, fmt.Errorf("invalid initialization value: %w", err)
				} else {
					mem.Bool[name] = val
				}
			case "integer":
				if val, err := getBasicMemoryInteger(initvalue); err != nil {
					return nil, fmt.Errorf("invalid initialization value: %w", err)
				} else {
					mem.Integer[name] = val
				}
			case "float":
				if val, err := getBasicMemoryFloat(initvalue); err != nil {
					return nil, fmt.Errorf("invalid initialization value: %w", err)
				} else {
					mem.Float[name] = val
				}
			case "text":
				if val, err := getBasicMemoryText(initvalue); err != nil {
					return nil, fmt.Errorf("invalid initialization value: %w", err)
				} else {
					mem.Text[name] = val
				}
			case "time":
				if val, err := getBasicMemoryTime(initvalue); err != nil {
					return nil, fmt.Errorf("invalid initialization value: %w", err)
				} else {
					mem.Time[name] = val
				}
			}
		}
	}
	return mem, nil
}

func getBasicMemoryBool(initvalue string) (bool, error) {
	if initvalue == "" {
		return false, nil
	} else {
		return strconv.ParseBool(initvalue)
	}
}

func getBasicMemoryInteger(initvalue string) (int64, error) {
	if initvalue == "" {
		return 0, nil
	} else {
		return strconv.ParseInt(initvalue, 10, 64)
	}
}

func getBasicMemoryFloat(initvalue string) (float64, error) {
	if initvalue == "" {
		return 0, nil
	} else {
		return strconv.ParseFloat(initvalue, 64)
	}
}

func getBasicMemoryText(initvalue string) (string, error) {
	return initvalue, nil
}

func getBasicMemoryTime(initvalue string) (time.Time, error) {
	if initvalue == "" {
		return time.Now(), nil
	} else {
		return time.Parse(time.Stamp, initvalue)
	}
}
