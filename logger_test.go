package niologrus

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-nio/nio"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TODO: add more unit tests
func TestMiddleware(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := logrus.New()
	logger.SetOutput(buf)
	logrusEntry := logrus.NewEntry(logger)

	e := nio.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := Middleware(logrusEntry)(func(c nio.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Status 2xx
	h(c)

	assert.Contains(t, buf.String(), "done")
}
