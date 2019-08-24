package handle

import (
	"context"
	"github.com/cgCodeLife/image_recognition_web/algorithm"
	"github.com/cgCodeLife/image_recognition_web/errno"
	"github.com/cgCodeLife/image_recognition_web/storage"
	"github.com/cgCodeLife/logs"
	"github.com/gin-gonic/gin"
	"time"
)

func GetFrameFaceInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(VID) == "" {
		return errno.InvalidVid
	}
	if c.Query(START_TIME) == "" || c.Query(END_TIME) == "" {
		return errno.InvalidTime
	}
	req := storage.GetFrameFaceInfoReq{Vid: c.Query(VID), StartTime: c.Query(START_TIME), EndTime: c.Query(END_TIME)}
	start := time.Now()
	frame, err := storage.GetFrameFaceInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameFaceInfo error=%s", err)
		if err == errno.INVALID_VID {
			return errno.InvalidVid
		}
		return errno.InternalErr
	}
	logs.CtxInfo(ctx, "GetFrameFaceInfo from storage record cost=%v", time.Since(start))

	//get frame now we need image recognition result
	start = time.Now()
	data, err := algorithm.GetFrameFaceInfo(frame)
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameFaceInfo error=%s", err)
		return errno.InternalErr
	}
	logs.CtxInfo(ctx, "GetFrameFaceInfo from algorithm cost=%v", time.Since(start))

	start = time.Now()
	res, ids, err := storage.GetFaceInfo(ctx, data)
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameFaceInfo error=%s", err)
		if err == errno.INVALID_BOUDINGBOX {
			return errno.InvalidBoudingBox
		}
		return errno.InternalErr
	}
	logs.CtxInfo(ctx, "GetFaceInfo from storage cost=%v", time.Since(start))

	//record for history search
	go storage.AddFrameFaceRecord(
		storage.FaceRecordInfo{
			Vid:      c.Query(VID),
			HumanIds: ids,
		})
	return errno.OK(res)
}

func GetFrameVehicleInfo(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(VID) == "" {
		return errno.InvalidVid
	}
	if c.Query(START_TIME) == "" || c.Query(END_TIME) == "" {
		return errno.InvalidTime
	}
	req := storage.GetFrameVehicleInfoReq{Vid: c.Query(VID), StartTime: c.Query(START_TIME), EndTime: c.Query(END_TIME)}
	frames, err := storage.GetFrameVehicleInfo(ctx, req)
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameVehicleInfo error=%s", err)
		if err == errno.INVALID_VID {
			return errno.InvalidVid
		}
		return errno.InternalErr
	}
	//get frame now we need image recognition result
	data, err := algorithm.GetFrameVehicleInfo(frames)
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameFaceInfo error=%s", err)
		return errno.InternalErr
	}
	//识别结果入库为历史查询做准备
	//todo.....
	//虽然图像识别识别出的信息有很多,但是这个接口只返回一个id及匹配到的车辆信息
	//need search db output vehicleinfo
	res := storage.GetFrameVehicleInfoRes{
		TestField: string(data),
	}
	return errno.OK(res)
}

func GetFrameVehicleTrafficFlow(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(VID) == "" {
		return errno.InvalidVid
	}
	if c.Query(START_TIME) == "" || c.Query(END_TIME) == "" {
		return errno.InvalidTime
	}
	req := storage.GetFrameVehicleInfoReq{Vid: c.Query(VID), StartTime: c.Query(START_TIME), EndTime: c.Query(END_TIME)}
	frames, err := storage.GetFrameVehicleInfo(ctx, req) //交通类的视频通用
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameVehicleTrafficFlow error=%s", err)
		if err == errno.INVALID_VID {
			return errno.InvalidVid
		}
		return errno.InternalErr
	}
	//get frame now we need image recognition result
	data, err := algorithm.GetFrameVehicleInfo(frames)
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameVehicleTrafficFlow error=%s", err)
		return errno.InternalErr
	}
	//识别结果入库为历史查询做准备
	//todo.....
	//虽然图像识别识别出的信息有很多,但是这个接口只返回一个交通流量的信息
	//need search db output vehicleinfo
	res := storage.GetFrameVehicleTrafficFlowRes{
		TestField: string(data),
	}
	return errno.OK(res)
}

func GetFrameVehicleAvgSpeed(c *gin.Context, ctx context.Context) errno.Payload {
	if c.Query(VID) == "" {
		return errno.InvalidVid
	}
	if c.Query(START_TIME) == "" || c.Query(END_TIME) == "" {
		return errno.InvalidTime
	}
	req := storage.GetFrameVehicleInfoReq{Vid: c.Query(VID), StartTime: c.Query(START_TIME), EndTime: c.Query(END_TIME)}
	frames, err := storage.GetFrameVehicleInfo(ctx, req) //交通类的视频通用
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameVehicleAvgSpeed error=%s", err)
		if err == errno.INVALID_VID {
			return errno.InvalidVid
		}
		return errno.InternalErr
	}
	//get frame now we need image recognition result
	data, err := algorithm.GetFrameVehicleInfo(frames)
	if err != nil {
		logs.CtxError(ctx, "method=GetFrameVehicleAvgSpeed error=%s", err)
		return errno.InternalErr
	}
	//识别结果入库为历史查询做准备
	//todo.....
	//虽然图像识别识别出的信息有很多,但是这个接口只返回一个平均车速的信息
	//need search db output vehicleinfo
	res := storage.GetFrameVehicleAvgSpeedRes{
		TestField: string(data),
	}
	return errno.OK(res)
}
