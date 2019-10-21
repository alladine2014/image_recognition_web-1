package storage

import (
	"context"
	"database/sql"
	"fmt"
	"work/image_recognition_web-1/algorithm"

	"work/image_recognition_web-1/errno"

	"work/image_recognition_web-1/util"

	"work/image_recognition_web-1/videolib"

	"github.com/cgCodeLife/logs"
)

type GetVehicleVideoInfoReq struct {
}

type GetVehicleVideoInfoRes struct {
	Vid        string       `json:"vid"`
	CameraInfo CameraInfoSt `json:"camera_info"`
}

type GetFaceVideoInfoReq struct {
	//need page default all
}

type GetFaceVideoInfoRes struct {
	Vid        string       `json:"vid"`
	CameraInfo CameraInfoSt `json:"camera_info"`
}

type CameraInfoSt struct {
	Uid      string     `json:"id"`
	Mac      string     `json:"mac"`
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
	Path string `json:"path"`
}

type GetFrameVehicleTrafficFlowReq struct {
	Vid       string
	StartTime string
	EndTime   string
}

type GetFrameVehicleTrafficFlowRes struct {
	TestField string
}

type GetFrameVehicleAvgSpeedReq struct {
	Vid       string
	StartTime string
	EndTime   string
}

type GetFrameVehicleAvgSpeedRes struct {
	TestField string
}

type GetFrameVehicleInfoReq struct {
	Vid       string
	StartTime string
	EndTime   string
}

type GetFrameVehicleInfoRes struct {
	//image recognition field
	TestField string
}

type GetFrameFaceInfoReq struct {
	Vid       string
	StartTime string
	EndTime   string
}

type GetFaceHistoryReq struct {
	HumanId   string
	StartTime string
	EndTime   string
}

type GetFaceHistoryRes struct {
	Vid        string       `json:"vid"`
	CameraInfo CameraInfoSt `json:"camera_info"`
	BaseInfo   BaseInfoSt   `json:"base_info"`
	Pic        string       `json:"pic"`
	// HumanId     string       `json:"id"` //human_id
	// Credibility float32      `json:"credibility"`
	Time string `json:"time"`
}

type AddFaceInfoReq struct {
	HumanId string `json:"id"`
	Pic     string `json:"pic"`
	Name    string `json:"name"`
	Time    string `json:"time"`
	Note    string `json:"note"`
}

type UpdateFaceInfoReq struct {
	HumanId string `json:"id"`
	Pic     string `json:"pic"`
	Name    string `json:"name"`
	Time    string `json:"time"`
	Note    string `json:"note"`
}

type UpdateFaceInfoRes struct {
	BaseInfo   BaseInfoSt   `json:"base_info"`
	CameraInfo CameraInfoSt `json:"camera_info"`
}

type SearchFaceInfoReq struct {
	Vid     string
	HumanId string
	Name    string
}

type SearchFaceInfoRes struct {
	BaseInfo   BaseInfoSt   `json:"base_info"`
	CameraInfo CameraInfoSt `json:"camera_info"`
}

type DeleteFaceInfoReq struct {
	HumanId string
}

type AddCameraInfoReq struct {
	Uid  string `json:"id"`
	Mac  string `json:"mac"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
	Addr string `json:"addr"`
	// Stream string `json:"stream"`
	// Vid    string `json:"vid"`
}

type UpdateCameraInfoReq struct {
	Uid  string `json:"id"`
	Mac  string `json:"mac"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
	Addr string `json:"addr"`
	// Stream string `json:"stream"`
	// Vid    string `json:"vid"`
}

type UpdateCameraInfoRes struct {
	Uid      string     `json:"id"`
	Mac      string     `json:"mac"`
	Location LocationSt `json:"location"`
	Addr     string     `json:"addr"`
	// Stream   string     `json:"stream"`
	// Vid      string     `json:"vid"`
}
type SearchCameraInfoReq struct {
	Uid string
}

// type DeleteCameraInfoReq struct {
// 	CameraID string
// }

type SearchCameraInfoRes struct {
	Uid      string     `json:"id"`
	Mac      string     `json:"mac"`
	Location LocationSt `json:"location"`
	Addr     string     `json:"addr"`
	// Stream   string     `json:"stream"`
	// Vid      string     `json:"vid"`
}

type AddVehicleInfoReq struct {
	Uid     string `json:"id"` //vehicle的id,车牌号码
	HumanId string `json:"human_id"`
	Time    string `json:"time"`
	Name    string `json:"name"`
	Note    string `json:"note"`
}

type UpdateVehicleInfoReq struct {
	Uid     string `json:"id"` //vehicle的id,车牌号码
	HumanId string `json:"human_id"`
	Time    string `json:"time"`
	Name    string `json:"name"`
	Note    string `json:"note"`
}

type UpdateVehicleInfoRes struct {
	Uid     string `json:"id"` //vehicle的id,车牌号码
	HumanId string `json:"human_id"`
	Time    string `json:"time"`
	Name    string `json:"name"`
	Note    string `json:"note"`
}

type SearchVehicleInfoReq struct {
	Uid string
}

type SearchVehicleInfoRes struct {
	Uid     string `json:"id"` //vehicle的id,车牌号码
	HumanId string `json:"human_id"`
	Time    string `json:"time"`
	Name    string `json:"name"`
	Note    string `json:"note"`
}

type GetFrameFaceInfoRes struct {
	Humans []HumanSt `json:"humans"`
}

type HumanSt struct {
	Location   BoudingBoxSt `json:"location"`
	BaseInfo   BaseInfoSt   `json:"base_info"`
	CameraInfo CameraInfoSt `json:"camera_info"`
}

type BoudingBoxSt struct {
	UpperLeftX   float64 `json:"upper_left_x"`
	UpperLeftY   float64 `json:"upper_left_y"`
	BottomRightX float64 `json:"bottom_right_x"`
	BottomRightY float64 `json:"bottom_right_y"`
}

type BaseInfoSt struct {
	Id   string `json:"id"`
	Pic  string `json:"pic"`
	Name string `json:"name"`
	Time string `json:"time"`
	Note string `json:"note"`
}

type FaceRecordInfo struct {
	Vid      string
	HumanIds []string
	Pic      string
}

func GetFaceInfo(ctx context.Context, faceInfo algorithm.FrameFaceRes) (GetFrameFaceInfoRes, []string, error) {
	res := GetFrameFaceInfoRes{}
	ids := make([]string, 0)
	res.Humans = make([]HumanSt, 0)
	for _, face := range faceInfo.Body.Face {
		name := face.Classification
		ress, err := SearchFaceInfo(ctx, SearchFaceInfoReq{Name: name})
		if err != nil {
			return res, ids, err
		}
		if len(face.BoudingBox) != 4 {
			return res, ids, errno.INVALID_BOUDINGBOX
		}
		if len(ress) == 0 {
			logs.Warnf("name=%s not found from db", name)
			continue
		}
		baseInfoRes := ress[0] //only one
		res.Humans = append(res.Humans, HumanSt{
			Location: BoudingBoxSt{
				UpperLeftX:   face.BoudingBox[0],
				UpperLeftY:   face.BoudingBox[1],
				BottomRightX: face.BoudingBox[2],
				BottomRightY: face.BoudingBox[3],
			},
			BaseInfo:   baseInfoRes.BaseInfo,
			CameraInfo: baseInfoRes.CameraInfo,
		})
		ids = append(ids, baseInfoRes.BaseInfo.Id)
	}
	return res, ids, nil
}

func AddFrameFaceRecord(faceRecordInfo FaceRecordInfo) {
	//get camera uid
	cameraUid, err := getCameraUidByVid(faceRecordInfo.Vid)
	if err != nil {
		logs.Errorf("faceRecordInfo:%+v error=%s", err)
		return
	}
	for _, id := range faceRecordInfo.HumanIds {
		sql := `insert into ` + humanFaceRecordTable
		sql += " (vid, human_id, time, credibility, uid, pic) values("
		sql += `"` + faceRecordInfo.Vid + `"`
		sql += " ," + `"` + id + `"`
		sql += " ," + `"` + util.GetTimeNow() + `"`
		sql += " ,1.00"
		sql += " ," + `"` + cameraUid + `"`
		sql += " ," + `"` + faceRecordInfo.Pic + `"`
		sql += ")"
		err = storage.dbQueryWrite(context.TODO(), sql)
		if err != nil {
			logs.Errorf("sql=%s error=%s", sql, err)
			continue
		}
	}
}

func getCameraUidByVid(vid string) (string, error) {
	sql := "select uid from " + videoTable
	sql += " where vid=" + `"` + vid + `"`
	uid, err := storage.dbQuery(context.TODO(), SEARCH_CAMERA_UID, sql)
	if err != nil {
		return "", err
	}
	return uid.(string), nil
}

func searchCameraUid(stm *sql.Stmt) (string, error) {
	var uid string
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("query error=%s", err)
		return uid, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&uid); err != nil {
			logs.Errorf("scan row error=%s", err)
			return uid, err
		}
	}
	return uid, nil
}

func UpdateVehicleInfo(ctx context.Context, req UpdateVehicleInfoReq, uid string) (UpdateVehicleInfoRes, error) {
	res := UpdateVehicleInfoRes{}
	// sql := fmt.Sprintf("update %s SET uid=%s,human_id=%s,time=%s,note=%s,name=%s where uid=%s)",
	// vehicleLibTable, req.HumanId, req.Time, req.Note, req.Name, req.Uid, uid)
	sql := `update ` + vehicleLibTable + ` SET`
	if req.Uid != "" {
		sql += " uid=" + req.Uid
	}
	if req.HumanId != "" {
		sql += " human_id=" + `"` + req.HumanId + `"` + ","
	}
	if req.Time != "" {
		sql += " time=" + `"` + req.Time + `"` + ","
	}
	if req.Note != "" {
		sql += " note=" + req.Note + ","
	}
	if req.Name != "" {
		sql += " name=" + req.Name + ","
	}
	if ',' == sql[len(sql)-1] {
		sql = sql[0 : len(sql)-1]
	}
	sql += " where uid=" + uid
	err := storage.dbQueryWrite(ctx, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return res, err
	}
	//select
	data, err := SearchVehicleInfo(ctx, SearchVehicleInfoReq{Uid: uid})
	if err != nil {
		return res, err
	}
	return UpdateVehicleInfoRes(data), nil
}

func SearchVehicleInfo(ctx context.Context, req SearchVehicleInfoReq) (SearchVehicleInfoRes, error) {
	res := SearchVehicleInfoRes{}
	sql := fmt.Sprintf("select human_id,time,note,name,uid from %s where uid=%s", vehicleLibTable, req.Uid)
	data, err := storage.dbQuery(ctx, SEARCH_VEHICLE_INFO, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return res, err
	}
	return data.(SearchVehicleInfoRes), nil
}

func searchVehicleInfo(stm *sql.Stmt) (SearchVehicleInfoRes, error) {
	res := SearchVehicleInfoRes{}
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("query error=%s", err)
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&res.HumanId, &res.Time, &res.Note, &res.Name, &res.Uid); err != nil {
			logs.Errorf("scan row error=%s", err)
			return res, err
		}
	}
	return res, nil
}

func AddVehicleInfo(ctx context.Context, req AddVehicleInfoReq) error {
	sql := fmt.Sprintf("insert into %s (human_id,time,note,name,uid) values(%s,%s,%s,%s,%s)",
		vehicleLibTable, req.HumanId, `"`+req.Time+`"`, `"`+req.Note+`"`, `"`+req.Name+`"`, req.Uid)
	err := storage.dbQueryWrite(ctx, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return err
	}
	return nil
}

func SearchFaceInfo(ctx context.Context, req SearchFaceInfoReq) ([]SearchFaceInfoRes, error) {
	var sql string
	var humanIdSearchCondition, nameSearchCondition, vidSearchCondition, searchCondition string
	if req.HumanId == "" {
		humanIdSearchCondition = ""
	} else if req.Name != "" || req.Vid != "" {
		humanIdSearchCondition = fmt.Sprintf("a.human_id=%s and", req.HumanId)
	} else {
		humanIdSearchCondition = fmt.Sprintf("a.human_id=%s", req.HumanId)
	}
	if req.Name == "" {
		nameSearchCondition = ""
	} else if req.Vid != "" {
		nameSearchCondition = fmt.Sprintf("a.name=%s and", `"`+req.Name+`"`)
	} else {
		nameSearchCondition = fmt.Sprintf("a.name=%s", `"`+req.Name+`"`)
	}
	if req.Vid == "" {
		vidSearchCondition = ""
	} else {
		vidSearchCondition = fmt.Sprintf("b.vid=%s", `"`+req.Vid+`"`)
	}
	if req.HumanId == "" && req.Name == "" && req.Vid == "" { //default
		// searchCondition = "b.vid ='v1111'"
		searchCondition = ""
	} else {
		searchCondition = fmt.Sprintf("%s %s %s", humanIdSearchCondition, nameSearchCondition, vidSearchCondition)
	}
	sql = fmt.Sprintf("select a.human_id, a.pic, a.name, a.time, a.note, b.uid, b.mac, b.addr, b.lat, b.lon from %s as a, %s as b where %s", faceLibTable, cameraTable, searchCondition)
	data, err := storage.dbQuery(ctx, SEARCH_FACE_INFO, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return nil, err
	}
	return data.([]SearchFaceInfoRes), nil
}

func searchFaceInfo(stm *sql.Stmt) ([]SearchFaceInfoRes, error) {
	res := make([]SearchFaceInfoRes, 0)
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("query error=%s", err)
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		item := SearchFaceInfoRes{}
		baseInfo := BaseInfoSt{}
		cameraInfo := CameraInfoSt{}
		locationInfo := LocationSt{}
		if err := rows.Scan(&baseInfo.Id, &baseInfo.Pic, &baseInfo.Name, &baseInfo.Time, &baseInfo.Note, &cameraInfo.Uid, &cameraInfo.Mac, &cameraInfo.Addr, &locationInfo.Lat, &locationInfo.Lon); err != nil {
			logs.Errorf("scan row error=%s", err)
			return res, err
		}
		item.BaseInfo = baseInfo
		cameraInfo.Location = locationInfo
		item.CameraInfo = cameraInfo
		res = append(res, item)
	}
	return res, nil
}

func AddCameraInfo(ctx context.Context, req AddCameraInfoReq) error {
	sql := fmt.Sprintf("insert into %s (uid,mac,addr,lat,lon) values(%s,%s,%s,%s,%s)",
		cameraTable, `"`+req.Uid+`"`, `"`+req.Mac+`"`, `"`+req.Addr+`"`, req.Lat, req.Lon)
	// sql := fmt.Sprintf("insert into %s (uid,mac,addr,lat,lon,stream,vid) values(%s,%s,%s,%s,%s,%s,%s)",
	// 	cameraTable, `"`+req.Uid+`"`, `"`+req.Mac+`"`, `"`+req.Addr+`"`, req.Lat, req.Lon,
	// 	`"`+req.Stream+`"`, `"`+req.Vid+`"`)
	err := storage.dbQueryWrite(ctx, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return err
	}
	return nil
}

func UpdateCameraInfo(ctx context.Context, req UpdateCameraInfoReq, uid string) (UpdateCameraInfoRes, error) {
	res := UpdateCameraInfoRes{}
	// sql := fmt.Sprintf("update %s SET uid=%s,mac=%s,addr=%s,lat=%f, lon=%f where uid=%s)",
	// cameraTable, req.Uid, req.Mac, req.Addr, req.Lat, req.Lon, uid)
	sql := `update ` + cameraTable + ` SET`
	if req.Uid != "" {
		sql += " uid=" + req.Uid
		// sql += " uid=" + `"` + req.Uid + `"` + ","
	}
	if req.Mac != "" {
		sql += " mac=" + `"` + req.Mac + `"` + ","
	}
	if req.Addr != "" {
		sql += " addr=" + `"` + req.Addr + `"` + ","
	}
	if req.Lat != "" {
		sql += " lat=" + req.Lat + ","
	}
	if req.Lon != "" {
		sql += " lon=" + req.Lon + ","
	}
	// if req.Stream != "" {
	// 	sql += " stream=" + `"` + req.Stream + `"` + ","
	// }
	// if req.Vid != "" {
	// 	sql += " vid=" + `"` + req.Vid + `"` + ","
	// }
	if ',' == sql[len(sql)-1] {
		sql = sql[0 : len(sql)-1]
	}
	sql += " where uid=" + uid
	// sql += " where uid='" + uid + `'`
	// logs.CtxError(ctx, "albin debug req=%s ", req)
	// logs.CtxError(ctx, "albin debug sql=%s ", sql)
	err := storage.dbQueryWrite(ctx, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return res, err
	}
	//select
	data, err := SearchCameraInfo(ctx, SearchCameraInfoReq{Uid: uid})
	if err != nil {
		return res, err
	}
	return UpdateCameraInfoRes(data[0]), nil
}

func SearchCameraInfo(ctx context.Context, req SearchCameraInfoReq) ([]SearchCameraInfoRes, error) {
	res := make([]SearchCameraInfoRes, 0)
	var sql string
	if req.Uid != "" {
		sql = fmt.Sprintf("select uid,mac,addr,lat,lon from %s where uid=%s", cameraTable, req.Uid)
		// sql = fmt.Sprintf("select uid,mac,addr,lat,lon,stream,vid from %s where uid='%s'", cameraTable, req.Uid)
	} else {
		sql = fmt.Sprintf("select uid,mac,addr,lat,lon from %s", cameraTable)
		// sql = fmt.Sprintf("select uid,mac,addr,lat,lon,stream,vid from %s", cameraTable)
	}
	data, err := storage.dbQuery(ctx, SEARCH_CAMERA_INFO, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return res, err
	}
	return data.([]SearchCameraInfoRes), nil
}

func searchCameraInfo(stm *sql.Stmt) ([]SearchCameraInfoRes, error) {
	res := make([]SearchCameraInfoRes, 0)
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("query error=%s", err)
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		item := SearchCameraInfoRes{}
		if err := rows.Scan(&item.Uid, &item.Mac, &item.Addr, &item.Location.Lat, &item.Location.Lon); err != nil {
			// if err := rows.Scan(&item.Uid, &item.Mac, &item.Addr, &item.Location.Lat, &item.Location.Lon, &item.Stream, &item.Vid); err != nil {
			logs.Errorf("scan row error=%s", err)
			return res, err
		}
		res = append(res, item)
	}
	return res, nil
}

func AddFaceInfo(ctx context.Context, req AddFaceInfoReq) error {
	sql := fmt.Sprintf("insert into %s (human_id,pic,name,time, note) values(%s,%s,%s,%s,%s)",
		faceLibTable, req.HumanId, `"`+req.Pic+`"`, `"`+req.Name+`"`, `"`+req.Time+`"`, `"`+req.Note+`"`)
	err := storage.dbQueryWrite(ctx, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return err
	}
	return nil
}

func UpdateFaceInfo(ctx context.Context, req UpdateFaceInfoReq, humanId string) (UpdateFaceInfoRes, error) {
	res := UpdateFaceInfoRes{}
	// sql := fmt.Sprintf("update %s SET human_id=%s,pic=%s,name=%s,time=%s, note=%s where human_id=%s)",
	// faceLibTable, req.HumanId, req.Pic, req.Name, req.Time, req.Note, humanId)
	sql := `update ` + faceLibTable + ` SET`
	if req.HumanId != "" {
		sql += " human_id=" + req.HumanId
	}
	if req.Pic != "" {
		sql += " pic=" + `"` + req.Pic + `"` + ","
	}
	if req.Name != "" {
		sql += " name=" + `"` + req.Name + `"` + ","
	}
	if req.Time != "" {
		sql += " time=" + `"` + req.Time + `"` + ","
	}
	if req.Note != "" {
		sql += " note=" + `"` + req.Note + `"`
	}
	if ',' == sql[len(sql)-1] {
		sql = sql[0 : len(sql)-1]
	}
	sql += " where human_id=" + humanId
	err := storage.dbQueryWrite(ctx, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return res, err
	}
	//select
	if req.HumanId != "" {
		humanId = req.HumanId
	}
	data, err := SearchFaceInfo(ctx, SearchFaceInfoReq{HumanId: humanId})
	if err != nil {
		return res, err
	}
	return UpdateFaceInfoRes(data[0]), nil
}

func DelFaceInfo(ctx context.Context, req DeleteFaceInfoReq) error {
	sql := fmt.Sprintf("delete from %s where human_id=%s",
		faceLibTable, req.HumanId)
	err := storage.dbQueryWrite(ctx, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return err
	}
	return nil
}

// func DelCameraInfo(ctx context.Context, req DeleteCameraInfoReq) error {
// 	sql := fmt.Sprintf("delete from %s where uid=%s",
// 		cameraTable, req.CameraID)
// 	err := storage.dbQueryWrite(ctx, sql)
// 	if err != nil {
// 		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
// 		return err
// 	}
// 	return nil
// }

func GetFrameVehicleInfo(ctx context.Context, req GetFrameVehicleInfoReq) (*videolib.Frame, error) {
	data := storage.SearchFrame(req.Vid, req.StartTime, req.EndTime)
	if data == nil {
		var err error
		//get file by vid
		file := storage.GetVideoFile(req.Vid)
		if file == "" {
			file, err = storage.GetUpdateVideoFile(req.Vid)
			if err != nil {
				logs.CtxError(ctx, "method=GetFrameVehicleInfo error=%s", err)
				return nil, err
			}
			if file == "" {
				logs.CtxError(ctx, "method=GetFrameVehicleInfo not found vid=%s", req.Vid)
				return nil, errno.INVALID_VID
			}
		}

		data, err = videolib.GetFrame(file, req.StartTime, req.EndTime)
		if err != nil {
			logs.CtxError(ctx, "method=GetFrameVehicleInfo error=%s", err)
			return nil, err
		}
		//update cache
		storage.CacheFrame(req.Vid, req.StartTime, req.EndTime, data)
	}
	return data, nil
}

func GetFrameFaceInfo(ctx context.Context, req GetFrameFaceInfoReq) (*videolib.Frame, error) {
	data := storage.SearchFrame(req.Vid, req.StartTime, req.EndTime)
	if data == nil {
		logs.CtxInfo(ctx, "star_time=%s end_time=%s not found in cache", req.StartTime, req.EndTime)
		var err error
		//get file by vid
		file := storage.GetVideoFile(req.Vid)
		if file == "" {
			logs.Infof("file is empty")
			file, err = storage.GetUpdateVideoFile(req.Vid)
			if err != nil {
				logs.CtxError(ctx, "method=GetFrameFaceInfo error=%s", err)
				return nil, err
			}
			if file == "" {
				logs.CtxError(ctx, "method=GetFrameFaceInfo not found vid=%s", req.Vid)
				return nil, errno.INVALID_VID
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

func GetVehicleVideoInfo(ctx context.Context, req GetVehicleVideoInfoReq) ([]GetVehicleVideoInfoRes, error) {
	//select from cache
	//todo....
	//select from db
	sql := fmt.Sprintf("select a.vid, a.uid, b.addr, b.lat, b.lon from %s as a, %s as b where a.uid=b.uid and type=1", videoTable, cameraTable)
	data, err := storage.dbQuery(ctx, GET_VEHICLE_VIDEO_INFO, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return nil, err
	}
	return data.([]GetVehicleVideoInfoRes), nil
}

func GetFaceHistory(ctx context.Context, req GetFaceHistoryReq) ([]GetFaceHistoryRes, error) {
	//select from cache
	//todo...
	//select from db
	sql := fmt.Sprintf("select a.vid, a.time, a.pic, b.uid, b.addr, b.lat, b.lon, c.human_id, c.pic, c.time, c.note, c.name from %s as a, %s as b, %s as c where a.uid=b.uid and a.human_id=c.human_id and a.time>= %s and a.time<= %s ", humanFaceRecordTable, cameraTable, faceLibTable, `"`+req.StartTime+`"`, `"`+req.EndTime+`"`)
	data, err := storage.dbQuery(ctx, GET_FACE_HISTORY, sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s error=%s", sql, err)
		return nil, err
	}
	return data.([]GetFaceHistoryRes), nil
}

func getFaceHistory(stm *sql.Stmt) ([]GetFaceHistoryRes, error) {
	res := make([]GetFaceHistoryRes, 0)
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("query error=%s", err)
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		item := GetFaceHistoryRes{}
		if err := rows.Scan(&item.Vid, &item.Time, &item.Pic, &item.CameraInfo.Uid,
			&item.CameraInfo.Addr, &item.CameraInfo.Location.Lat,
			&item.CameraInfo.Location.Lon, &item.BaseInfo.Id, &item.BaseInfo.Pic, &item.BaseInfo.Time, &item.BaseInfo.Note, &item.BaseInfo.Name); err != nil {
			logs.Errorf("scan row error=%s", err)
			return res, err
		}
		res = append(res, item)
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
		if err := rows.Scan(&item.Vid, &item.CameraInfo.Uid,
			&item.CameraInfo.Addr, &item.CameraInfo.Location.Lat,
			&item.CameraInfo.Location.Lon); err != nil {
			logs.Errorf("scan row error=%s", err)
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func getVehicleVideoInfo(stm *sql.Stmt) ([]GetVehicleVideoInfoRes, error) {
	res := make([]GetVehicleVideoInfoRes, 0)
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("query error=%s", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := GetVehicleVideoInfoRes{}
		if err := rows.Scan(&item.Vid, &item.CameraInfo.Uid,
			&item.CameraInfo.Addr, &item.CameraInfo.Location.Lat,
			&item.CameraInfo.Location.Lon); err != nil {
			logs.Errorf("scan row error=%s", err)
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}
