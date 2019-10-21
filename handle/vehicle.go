package handle

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"work/image_recognition_web-1/errno"

	"work/image_recognition_web-1/storage"

	"github.com/cgCodeLife/logs"
	"github.com/gin-gonic/gin"
)

func GetVehicleVideoInfo(c *gin.Context, ctx context.Context) errno.Payload {
	io.Copy(ioutil.Discard, c.Request.Body)
	//query storage
	req := storage.GetVehicleVideoInfoReq{}
	data, err := storage.GetVehicleVideoInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "methd=GetVehicleVideoInfo error=%s", err)
		return errno.SearchDBErr
	}

	return errno.OK(data)
}

func AddVehicleInfo(c *gin.Context, ctx context.Context) errno.Payload {
	req, err := parseAddVehicleInfoReq(c)
	if err != nil {
		return errno.ParseJsonError
	}
	err = storage.AddVehicleInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=AddVehicleInfo error=%s", err)
		return errno.AddDBErr
	}
	return errno.OK(nil)
}

func UpdateVechicleInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(VEHICLE_ID) == "" {
		return errno.InvalidVehicleId
	}
	req, err := parseUpdateVehicleInfoReq(c)
	if err != nil {
		return errno.ParseJsonError
	}
	data, err := storage.UpdateVehicleInfo(ctx, req, c.Query(VEHICLE_ID))
	if err != nil {
		logs.CtxError(ctx, "method=UpdateVehicleInfo error=%s", err)
		return errno.UpdateDBErr
	}
	return errno.OK(data)
}

func SearchVehicleInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(VEHICLE_ID) == "" {
		return errno.InvalidVehicleId
	}
	req := storage.SearchVehicleInfoReq{Uid: c.Query(VEHICLE_ID)}
	data, err := storage.SearchVehicleInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=SearchVehicleInfo error=%s", err)
		return errno.SearchDBErr
	}
	return errno.OK(data)
}

func parseAddVehicleInfoReq(c *gin.Context) (storage.AddVehicleInfoReq, error) {
	res := storage.AddVehicleInfoReq{}
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}
	return res, nil
}

func parseUpdateVehicleInfoReq(c *gin.Context) (storage.UpdateVehicleInfoReq, error) {
	res := storage.UpdateVehicleInfoReq{}
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
		case "human_id":
			res.HumanId = field.Value
		case DATE:
			res.Time = field.Value
		case HUMAN_NAME:
			res.Name = field.Value
		case NOTE:
			res.Note = field.Value
		case VEHICLE_ID:
			res.Uid = field.Value
		}
	}
	return res, nil
}
