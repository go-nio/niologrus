package niologrus

import (
	"github.com/go-nio/nio"
)

type (
	// Config defines the config for Logger middleware.
	Config struct {
		// Skipper defines a function to skip middleware.
		Skipper nio.Skipper
	}
)

var (
	// DefaulConfig is the default Logger middleware config.
	DefaulConfig = Config{}
)

// Middleware returns a middleware which logs requests
func Middleware() nio.MiddlewareFunc {
	return MiddlewareWithConfig(DefaulConfig)
}

// MiddlewareWithConfig returns a Logger middleware with config.
// See: `Middleware()`.
func MiddlewareWithConfig(config Config) nio.MiddlewareFunc {
	return func(next nio.HandlerFunc) nio.HandlerFunc {
		return func(c nio.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			return next(c)
		}
	}
}
