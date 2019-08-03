package handle

import (
	"context"
	"github.com/cgCodeLife/image_recognition_web/errno"
	"github.com/cgCodeLife/image_recognition_web/storage"
	"github.com/cgCodeLife/logs"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
)

func GetFaceVideoInfo(c *gin.Context, ctx context.Context) errno.Payload {
	io.Copy(ioutil.Discard, c.Request.Body) // discard body anyway
	//query storage
	req := &storage.GetFaceVideoInfoReq{}
	data, err := storage.GetFaceVideoInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "methd=GetFaceVideoInfo error=%s", err)
		return errno.SearchDBErr
	}

	return errno.OK(data)
}

func GetFaceHistory(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}

func AddFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}

func SearchFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}

func UpdateFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}

func DelFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	return errno.OK(nil)
}
