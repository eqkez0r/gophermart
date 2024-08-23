package e

import "fmt"

func Wrap(point string, err error) error {
	return fmt.Errorf("%s: %s", point, err.Error())
}
