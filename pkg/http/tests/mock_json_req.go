package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

func MockJSONRequest(c *gin.Context, method, cType string, content interface{}) {
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", cType)
	b, _ := json.Marshal(content)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
}
