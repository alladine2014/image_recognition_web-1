package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Options serves Options method api
func Options(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Headers", "X-TT-Access, Content-Type, accept, content-disposition, content-range")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.String(http.StatusOK, "")
}
