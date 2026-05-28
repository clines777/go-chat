package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type responseRecorder struct {
	gin.ResponseWriter
	buf    *bytes.Buffer
	status int
}

const respBufLimit = 1024 * 4 // 4KB

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.buf.Len() < respBufLimit {
		remaining := respBufLimit - r.buf.Len()
		if len(b) > remaining {
			r.buf.Write(b[:remaining])
		} else {
			r.buf.Write(b)
		}
	}
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var reqBody []byte
		ct := c.Request.Header.Get("Content-Type")
		if c.Request.Body != nil && !strings.HasPrefix(ct, "multipart/") {
			reqBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, respBufLimit))
			c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewBuffer(reqBody), c.Request.Body))
		}

		rec := &responseRecorder{
			ResponseWriter: c.Writer,
			buf:            &bytes.Buffer{},
			status:         http.StatusOK,
		}
		c.Writer = rec

		c.Next()

		log.Printf("[HTTP] %s %s | status=%d latency=%s | req=%s | resp=%s",
			c.Request.Method,
			c.Request.URL.Path,
			rec.status,
			time.Since(start).Round(time.Millisecond),
			truncate(reqBody, 200),
			truncate(rec.buf.Bytes(), 200),
		)
	}
}

func truncate(b []byte, n int) string {
	if len(b) == 0 {
		return "-"
	}
	if len(b) > n {
		return string(b[:n]) + "…"
	}
	return string(b)
}
