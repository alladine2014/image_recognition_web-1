package ginex

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Engine struct {
	*gin.Engine
}

func New() *Engine {
	r := &Engine{
		Engine: gin.New(),
	}
	return r
}

func SetMode(mode string) {
	gin.SetMode(mode)
}

// write to applog and gin.DefaultErrorWriter
type recoverWriter struct{}

func (rw *recoverWriter) Write(p []byte) (int, error) {
	return gin.DefaultErrorWriter.Write(p)
}

func Default() *Engine {
	r := New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.RecoveryWithWriter(&recoverWriter{}))
	return r
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It also starts a pprof debug server and report framework meta info to bagent
func (engine *Engine) Run(addr string) error {
	return engine.Engine.Run(addr)
}
