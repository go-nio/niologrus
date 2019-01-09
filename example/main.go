package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/go-nio/nio"
	"github.com/go-nio/niologrus"
)

func main() {
	// create logrus entry
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logrusEntry := logrus.NewEntry(logger)

	// set nio to use logrus as it's internal logger
	n := nio.New(nio.WithLogger(logrusEntry))

	// pass logrus middleware to nio
	n.Use(niologrus.Middleware(logrusEntry))

	n.GET("/", func(c nio.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	http.ListenAndServe(":9000", n)
}
