package digutil

import (
	"github.com/samber/lo"
	"go.uber.org/dig"
)

// Get instance from dig.Container
func Get[T any](c *dig.Container) (T, error) {
	var t T
	err := c.Invoke(func(tt T) {
		t = tt
	})
	return t, err
}

func MustGet[T any](c *dig.Container) T {
	return lo.Must(Get[T](c))
}
