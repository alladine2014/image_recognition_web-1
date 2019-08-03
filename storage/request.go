package storage

import (
	"context"
	"database/sql"
	"fmt"
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

func GetFaceVideoInfo(ctx context.Context, req *GetFaceVideoInfoReq) ([]GetFaceVideoInfoRes, error) {
	//select from cache
	//todo....
	//select from db
	sql := fmt.Sprintf("select a.vid, a.uid, b.addr, b.lat, b.lon from %s as a, %s as b where a.uid=b.uid and type=0", videoTable, cameraTable)
	data, err := storage.dbQuery(ctx, GET_FACE_VIDEO_INFO, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", err)
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
