package handle

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"work/image_recognition_web-1/errno"
	"work/image_recognition_web-1/storage"

	"github.com/cgCodeLife/logs"
	"github.com/gin-gonic/gin"
)

func AddCameraInfo(c *gin.Context, ctx context.Context) errno.Payload {
	// //参数校验
	// if c.Query(CAMERA_ID) == "" {
	// 	return errno.InvalidCameraId
	// }
	// if c.Query(MAC) == "" {
	// 	return errno.InvalidMac
	// }
	// if c.Query(ADDR) == "" {
	// 	return errno.InvalidAddr
	// }
	// if c.Query(LAT) == "" || c.Query(LON) == "" {
	// 	return errno.InvalidLocation
	// }
	req, err := parseAddCameraInfoReq(c)
	if err != nil {
		return errno.ParseJsonError
	}
	// req := storage.AddCameraInfoReq{Uid: c.Query(CAMERA_ID), Mac: c.Query(MAC), Addr: c.Query(ADDR), Lat: c.Query(LAT), Lon: c.Query(LON)}
	err = storage.AddCameraInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=AddCameraInfo error=%s", err)
		return errno.AddDBErr
	}
	return errno.OK(nil)
}

func UpdateCameraInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(CAMERA_ID) == "" {
		return errno.InvalidCameraId
	}
	req, err := parseUpdateCameraInfoReq(c)
	if err != nil {
		return errno.ParseJsonError
	}
	data, err := storage.UpdateCameraInfo(ctx, req, c.Query(CAMERA_ID))
	if err != nil {
		logs.CtxError(ctx, "method=UpdateCameraInfo error=%s", err)
		return errno.UpdateDBErr
	}
	return errno.OK(data)
}

func parseUpdateCameraInfoReq(c *gin.Context) (storage.UpdateCameraInfoReq, error) {
	res := storage.UpdateCameraInfoReq{}
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
		case CAMERA_ID:
			res.Uid = field.Value
		case MAC:
			res.Mac = field.Value
		case ADDR:
			res.Addr = field.Value
		case LAT:
			res.Lat = field.Value
		case LON:
			res.Lon = field.Value
			// case STREAM:
			// 	res.Stream = field.Value
			// case VID:
			// 	res.Vid = field.Value
		}
	}
	return res, nil
}

func parseAddCameraInfoReq(c *gin.Context) (storage.AddCameraInfoReq, error) {
	res := storage.AddCameraInfoReq{}
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}
	return res, nil
}

func SearchCameraInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(CAMERA_ID) == "" {
		logs.CtxInfo(ctx, "like search")
	}
	req := storage.SearchCameraInfoReq{Uid: c.Query(CAMERA_ID)}
	data, err := storage.SearchCameraInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=SearchCameraInfo error=%s", err)
		return errno.SearchDBErr
	}
	return errno.OK(data)
}

// func DelCameraInfo(c *gin.Context, ctx context.Context) errno.Payload {
// 	if c.Query(CAMERA_ID) == "" {
// 		return errno.InvalidCameraId
// 	}
// 	req := storage.DeleteCameraInfoReq{CameraID: c.Query(CAMERA_ID)}
// 	err := storage.DelCameraInfo(ctx, req)
// 	if err != nil {
// 		logs.CtxError(ctx, "method=DeleteCameraInfo error=%s", err)
// 		if err == errno.NOT_FOUND {
// 			return errno.NotFoundRecord
// 		}
// 		return errno.DelDBErr
// 	}
// 	return errno.OK(nil)
// }
