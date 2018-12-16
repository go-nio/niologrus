package niologrus

import (
	"strconv"
	"time"

	"github.com/go-nio/nio"
	"github.com/sirupsen/logrus"
)

type options struct {
	skipper nio.Skipper
}

var defaultOptions = options{
	skipper: nio.DefaultSkipper,
}

// Option is middleware option
type Option func(*options)

// WithSkipper allows to pass middleware skipper
func WithSkipper(skipper nio.Skipper) Option {
	return func(o *options) {
		o.skipper = skipper
	}
}

// Middleware returns a middleware which logs requests
func Middleware(logrusEntry *logrus.Entry, opt ...Option) nio.MiddlewareFunc {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	return func(next nio.HandlerFunc) nio.HandlerFunc {
		return func(c nio.Context) error {
			if opts.skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			stop := time.Now()

			// get content length
			cl := req.Header.Get(nio.HeaderContentLength)
			if cl == "" {
				cl = "0"
			}

			// create default logger fields
			infoFields := logrus.Fields{
				"bytes_in":   cl,
				"bytes_out":  strconv.FormatInt(res.Size, 10),
				"host":       req.Host,
				"start_time": start.Format(time.RFC3339),
				"end_time":   stop.Format(time.RFC3339),
				"status":     res.Status,
				"method":     req.Method,
				"path":       req.URL.Path,
				"latency_ns": strconv.FormatInt(int64(stop.Sub(start)), 10),
				"latency":    stop.Sub(start).String(),
			}
			entry := logrusEntry.WithFields(infoFields)

			// log message
			if err != nil && res.Status >= 500 {
				entry.Error(err)
			} else {
				entry.Info("done")
			}

			return nil
		}
	}
}
