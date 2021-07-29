package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HasPostFormEmpty(c *gin.Context, keys ...string) string {
	for _, v := range keys {
		if strings.TrimSpace(c.PostForm(v)) == "" {
			return v
		}
	}
	return ""
}

func HasQueryEmpty(c *gin.Context, keys ...string) string {
	for _, v := range keys {
		if strings.TrimSpace(c.Query(v)) == "" {
			return v
		}
	}
	return ""
}

func ErrorResponse(c *gin.Context, code int, msg string, err error) {
	errorMsg := msg
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %v", msg, err)
	}
	c.Set("ErrorMsg", errorMsg)
	c.JSON(code, gin.H{"error": msg})
}

func APILogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		c.Next()

		header := ""
		for k, v := range c.Request.Header {
			header += k + ": " + fmt.Sprint(v) + "\r\n"
		}

		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path += "?" + c.Request.URL.RawQuery
		}

		msg := "------------------------------------------------------------\r\n"
		msg += fmt.Sprintf("[%s] | %s\r\n\r\n", time.Now().Format("2006/01/02 15:04:05"), c.ClientIP())
		msg += fmt.Sprintf("[Request] \r\n%s %s\r\n\r\n", c.Request.Method, path)
		msg += fmt.Sprintf("[Header] \r\n%s\r\n", header)
		msg += fmt.Sprintf("[Body] \r\n%s\r\n\r\n", string(body))
		msg += fmt.Sprintf("[Status] \r\n%v\r\n\r\n", c.Writer.Status())
		msg += fmt.Sprintf("[Response Data] \r\n%s\r\n\r\n", strings.TrimRight(blw.body.String(), "\n"))
		msg += fmt.Sprintf("[ErrorMsg] \r\n%v\r\n\r\n", c.Keys["ErrorMsg"])
		msg += "------------------------------------------------------------\r\n"

		go WriteLog("", msg)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
