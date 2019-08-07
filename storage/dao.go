package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cgCodeLife/image_recognition_web/videolib"
	"github.com/cgCodeLife/logs"
)

type GetFaceVideoInfoReq struct {
	//need page default all
}

type GetFaceVideoInfoRes struct {
	Vid        string       `json:"vid"`
	CameraInfo CameraInfoSt `json:"camera_info"`
}

type CameraInfoSt struct {
	Id       string     `json:"id"`
	Addr     string     `json:"addr"`
	Location LocationSt `json:"location"`
}

type LocationSt struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type GetVideoInfoReq struct {
	Vid string
}

type GetVideoInfoRes struct {
	Path string
}

type GetFrameFaceInfoReq struct {
	Vid       string
	StartTime string
	EndTime   string
}

type GetFrameFaceInfoRes struct {
	//image recognition field
	TestField string
}

func GetFrameFaceInfo(ctx context.Context, req GetFrameFaceInfoReq) ([]videolib.Frame, error) {
	data := storage.SearchFrame(req.Vid, req.StartTime, req.EndTime)
	if data == nil {
		var err error
		//get file by vid
		file := storage.GetVideoFile(req.Vid)
		if file == "" {
			file, err = storage.GetUpdateVideoFile(req.Vid)
			if err != nil {
				logs.CtxError(ctx, "method=GetFrameFaceInfo error=%s", err)
				return nil, err
			}
		}

		data, err = videolib.GetFrame(file, req.StartTime, req.EndTime)
		if err != nil {
			logs.CtxError(ctx, "method=GetFrameFaceInfo error=%s", err)
			return nil, err
		}
		//update cache
		storage.CacheFrame(req.Vid, req.StartTime, req.EndTime, data)
	}
	return data, nil
}

func GetVideoInfo(ctx context.Context, req GetVideoInfoReq) (GetVideoInfoRes, error) {
	//get cache
	//todo....
	sql := fmt.Sprintf("select path from %s where vid=\"%s\"", videoTable, req.Vid)
	data, err := storage.dbQuery(ctx, GET_VIDEO_INFO, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return GetVideoInfoRes{}, err
	}
	return data.(GetVideoInfoRes), nil
}

func getVideoInfo(stm *sql.Stmt) (GetVideoInfoRes, error) {
	res := GetVideoInfoRes{}
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("query error=%s", err)
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&res.Path); err != nil {
			logs.Errorf("scan row error=%s", err)
			return res, err
		}
	}
	return res, nil
}

func GetFaceVideoInfo(ctx context.Context, req GetFaceVideoInfoReq) ([]GetFaceVideoInfoRes, error) {
	//select from cache
	//todo....
	//select from db
	sql := fmt.Sprintf("select a.vid, a.uid, b.addr, b.lat, b.lon from %s as a, %s as b where a.uid=b.uid and type=0", videoTable, cameraTable)
	data, err := storage.dbQuery(ctx, GET_FACE_VIDEO_INFO, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return nil, err
	}
	return data.([]GetFaceVideoInfoRes), nil
}

func getFaceVideoInfo(stm *sql.Stmt) ([]GetFaceVideoInfoRes, error) {
	res := make([]GetFaceVideoInfoRes, 0)
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("query error=%s", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := GetFaceVideoInfoRes{}
		if err := rows.Scan(&item.Vid, &item.CameraInfo.Id,
			&item.CameraInfo.Addr, &item.CameraInfo.Location.Lat,
			&item.CameraInfo.Location.Lon); err != nil {
			logs.Errorf("scan row error=%s", err)
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}
