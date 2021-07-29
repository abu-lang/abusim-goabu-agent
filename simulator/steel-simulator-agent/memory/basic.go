package memory

import (
	"fmt"
	"steel-lang/datastructure"
	"strconv"
	"time"
)

// New creates a new memory, based on the basic resources
func NewBasicMemory(items map[string]map[string][]string) (datastructure.ResourceController, error) {
	// I create an empty basic memory...
	mem := datastructure.MakeResources()
	// ... and I range over the items to initialize it
	for vartype, values := range items {
		for name, initvalues := range values {
			switch vartype {
			case "bool":
				val, err := assingBool(initvalues)
				if err != nil {
					return nil, err
				}
				mem.Bool[name] = *val
			case "integer":
				val, err := assingInteger(initvalues)
				if err != nil {
					return nil, err
				}
				mem.Integer[name] = *val
			case "float":
				val, err := assingFloat(initvalues)
				if err != nil {
					return nil, err
				}
				mem.Float[name] = *val
			case "text":
				val, err := assingText(initvalues)
				if err != nil {
					return nil, err
				}
				mem.Text[name] = *val
			case "time":
				val, err := assingTime(initvalues)
				if err != nil {
					return nil, err
				}
				mem.Time[name] = *val
			default:
				return nil, fmt.Errorf("unknown type \"%s\"", vartype)
			}
		}
	}
	return mem, nil
}

func assingBool(initvalues []string) (*bool, error) {
	switch len(initvalues) {
	case 0:
		def := false
		return &def, nil
	case 1:
		val, err := strconv.ParseBool(initvalues[0])
		if err != nil {
			return nil, err
		}
		return &val, nil
	default:
		return nil, ErrInvalidValues
	}
}

func assingInteger(initvalues []string) (*int64, error) {
	switch len(initvalues) {
	case 0:
		def := int64(0)
		return &def, nil
	case 1:
		val, err := strconv.ParseInt(initvalues[0], 10, 64)
		if err != nil {
			return nil, err
		}
		return &val, nil
	default:
		return nil, ErrInvalidValues
	}
}

func assingFloat(initvalues []string) (*float64, error) {
	switch len(initvalues) {
	case 0:
		def := float64(0)
		return &def, nil
	case 1:
		val, err := strconv.ParseFloat(initvalues[0], 64)
		if err != nil {
			return nil, err
		}
		return &val, nil
	default:
		return nil, ErrInvalidValues
	}
}

func assingText(initvalues []string) (*string, error) {
	switch len(initvalues) {
	case 0:
		def := ""
		return &def, nil
	case 1:
		return &initvalues[0], nil
	default:
		return nil, ErrInvalidValues
	}
}

func assingTime(initvalues []string) (*time.Time, error) {
	switch len(initvalues) {
	case 0:
		def := time.Now()
		return &def, nil
	case 1:
		val, err := time.Parse(time.Stamp, initvalues[0])
		if err != nil {
			return nil, err
		}
		return &val, nil
	default:
		return nil, ErrInvalidValues
	}
}
