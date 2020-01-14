package apigroup

import (
	"errors"
)

// SetAPIGroup function will return the api prefix for echo group
func SetAPIGroup(serviceType, serviceVersion string) string {
	var prefix string

	switch serviceType {
	case "backend-golang":
		prefix = "/api/" + serviceVersion
		break
	default:
		panic(errors.New("Service type does not support"))
	}

	return prefix
}
