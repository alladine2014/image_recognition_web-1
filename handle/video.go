package handle

import (
	"context"
	"github.com/cgCodeLife/image_recognition_web/errno"
	"github.com/cgCodeLife/image_recognition_web/storage"
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
