// Package ginbunyan provides log handling using bunyan package.
// Code structure based on zip package.
package ginbunyan

import (
	"time"

	"github.com/bhoriuchi/go-bunyan/bunyan"
	"github.com/gin-gonic/gin"
)

// logEntry represents a log line
type logEntry struct {
	status    int
	method    string
	path      string
	query     string
	ip        string
	userAgent string
	time      string
	latency   time.Duration
}

// Ginbunyan returns a gin.HandlerFunc (middleware) that logs requests using bhoriuchi/go-bunyan.
func Ginbunyan(logger *bunyan.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			le := logEntry{
				status:    c.Writer.Status(),
				method:    c.Request.Method,
				path:      path,
				query:     query,
				ip:        c.ClientIP(),
				userAgent: c.Request.UserAgent(),
				time:      end.Format(time.RFC3339),
				latency:   latency,
			}
			logger.Info(le)
		}
	}
}
