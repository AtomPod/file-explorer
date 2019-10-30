package services

import "fmt"

func invalidArgument(service, arg, method string) error {
	return fmt.Errorf("%s: invalid argument '%s' in %s", service, arg, method)
}