package memory

import (
	"steel-lang/datastructure"
	"strconv"
)

func NewBasicMemory(items map[string]map[string][]string) (datastructure.ResourceController, error) {
	mem := datastructure.MakeResources()
	for vartype, values := range items {
		for name, initvalues := range values {
			switch vartype {
			case "bool":
				val, err := assingBool(initvalues)
				if err != nil {
					return nil, err
				}
				mem.Bool[name] = val
			}
		}
	}
	return mem, nil
}

func assingBool(initvalues []string) (bool, error) {
	switch len(initvalues) {
	case 0:
		return false, nil
	case 1:
		val, err := strconv.ParseBool(initvalues[0])
		if err != nil {
			return false, err
		}
		return val, nil
	default:
		return false, ErrInvalidValues
	}
}
