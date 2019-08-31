package handle

import (
	"context"
	"encoding/json"
	"github.com/cgCodeLife/image_recognition_web/errno"
	"github.com/cgCodeLife/image_recognition_web/storage"
	"github.com/cgCodeLife/logs"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
)

type QueryDataSt struct {
	Data []DataSt `json:"data"`
}

type DataSt struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func GetFaceVideoInfo(c *gin.Context, ctx context.Context) errno.Payload {
	io.Copy(ioutil.Discard, c.Request.Body) // discard body anyway
	//query storage
	req := storage.GetFaceVideoInfoReq{}
	data, err := storage.GetFaceVideoInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "methd=GetFaceVideoInfo error=%s", err)
		return errno.SearchDBErr
	}

	return errno.OK(data)
}

func GetFaceHistory(c *gin.Context, ctx context.Context) errno.Payload {
	//参数校验
	if c.Query(FACE_ID) == "" {
		return errno.InvalidFaceId
	}
	if c.Query(START_TIME) == "" || c.Query(END_TIME) == "" {
		return errno.InvalidTime
	}
	req := storage.GetFaceHistoryReq{HumanId: c.Query(FACE_ID), StartTime: c.Query(START_TIME), EndTime: c.Query(END_TIME)}
	data, err := storage.GetFaceHistory(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "methd=GetFaceHistory error=%s", err)
		return errno.SearchDBErr
	}

	return errno.OK(data)
}

func AddFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	req, err := parseAddFaceInfoReq(c)
	if err != nil {
		return errno.ParseJsonError
	}
	err = storage.AddFaceInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=AddFaceInfo error=%s", err)
		return errno.AddDBErr
	}
	return errno.OK(nil)
}

func SearchFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	//参数校验
	if c.Query(FACE_ID) == "" && c.Query(HUMAN_NAME) == "" {
		logs.CtxInfo(ctx, "like serch")
	}
	if c.Query(VID) == "" {
		logs.CtxWarn(ctx, "vid null")
	}
	req := storage.SearchFaceInfoReq{HumanId: c.Query(FACE_ID), Name: c.Query(HUMAN_NAME)}
	data, err := storage.SearchFaceInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=SearchFaceInfo error=%s", err)
		return errno.SearchDBErr
	}
	return errno.OK(data)
}

func UpdateFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(FACE_ID) == "" {
		return errno.InvalidFaceId
	}
	req, err := parseUpdateFaceInfoReq(c)
	if err != nil {
		return errno.ParseJsonError
	}
	data, err := storage.UpdateFaceInfo(ctx, req, c.Query(FACE_ID))
	if err != nil {
		logs.CtxError(ctx, "method=UpdateFaceInfo error=%s", err)
		return errno.UpdateDBErr
	}
	return errno.OK(data)
}

func DelFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(FACE_ID) == "" {
		return errno.InvalidFaceId
	}
	req := storage.DeleteFaceInfoReq{HumanId: c.Query(FACE_ID)}
	err := storage.DelFaceInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=DeleteFaceInfo error=%s", err)
		if err == errno.NOT_FOUND {
			return errno.NotFoundRecord
		}
		return errno.DelDBErr
	}
	return errno.OK(nil)
}

func parseUpdateFaceInfoReq(c *gin.Context) (storage.UpdateFaceInfoReq, error) {
	res := storage.UpdateFaceInfoReq{}
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return res, err
	}
	var v QueryDataSt
	v.Data = make([]DataSt, 0)
	if err := json.Unmarshal(data, &v); err != nil {
		return res, err
	}
	for _, field := range v.Data {
		switch field.Key {
		case FACE_ID:
			res.HumanId = field.Value
		case PIC:
			res.Pic = field.Value
		case HUMAN_NAME:
			res.Name = field.Value
		case DATE:
			res.Time = field.Value
		case NOTE:
			res.Note = field.Value
		}
	}
	return res, nil
}

func parseAddFaceInfoReq(c *gin.Context) (storage.AddFaceInfoReq, error) {
	res := storage.AddFaceInfoReq{}
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}
	return res, nil
}
