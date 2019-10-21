package handle

import (
	"context"
	"work/image_recognition_web-1/errno"

	"work/image_recognition_web-1/storage"

	"github.com/cgCodeLife/logs"
	"github.com/gin-gonic/gin"
)

func GetVideoInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(VID) == "" {
		return errno.InvalidVid
	}
	req := storage.GetVideoInfoReq{Vid: c.Query(VID)}
	res, err := storage.GetVideoInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=GetVideoInfo error=%s", err)
		return errno.SearchDBErr
	}

	//res is a local filename
	return errno.LocalStream(res.Path)
}
