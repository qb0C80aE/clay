package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}

// PreloadBody preloads body contents before binding, and put that into the Context
// It allows Binders to rebind from the same stream
func PreloadBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := c.Request.Body
		defer body.Close()

		byteArray, err := ioutil.ReadAll(body)
		if err != nil {
			panic(err)
		}

		c.Set("PreloadedBody", byteArray)
		c.Request.Body = nopCloser{bytes.NewBuffer(byteArray)}

		c.Next()
	}
}
