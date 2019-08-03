package handle

import (
	"context"
	"github.com/cgCodeLife/image_recognition_web/errno"
	"github.com/gin-gonic/gin"
)

func GetFrameFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}

func GetFrameVehicleInfo(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}

func GetFrameVehicleTrafficFlow(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}

func GetFrameVehicleAvgSpeed(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}
