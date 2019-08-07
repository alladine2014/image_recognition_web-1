package handle

import (
	"context"
	"github.com/cgCodeLife/image_recognition_web/algorithm"
	"github.com/cgCodeLife/image_recognition_web/errno"
	"github.com/cgCodeLife/image_recognition_web/storage"
	"github.com/cgCodeLife/logs"
	"github.com/gin-gonic/gin"
)

const ()

func GetFrameFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(VID) == "" {
		return errno.InvalidVid
	}
	if c.Query(START_TIME) == "" || c.Query(END_TIME) == "" {
		return errno.InvalidTime
	}
	req := storage.GetFrameFaceInfoReq{Vid: c.Query(VID), StartTime: c.Query(START_TIME), EndTime: c.Query(END_TIME)}
	frames, err := storage.GetFrameFaceInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameFaceInfo error=%s", err)
		return errno.InternalErr
	}
	//get frame now we need image recognition result
	data, err := algorithm.GetFrameInfo(frames)
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameFaceInfo error=%s", err)
	}
	res := storage.GetFrameFaceInfoRes{
		TestField: string(data),
	}
	return errno.OK(res)
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
