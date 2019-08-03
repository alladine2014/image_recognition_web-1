package handle

import (
	"context"
	"github.com/cgCodeLife/image_recognition_web/errno"
	"github.com/gin-gonic/gin"
)

func GetVideoInfo(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}
