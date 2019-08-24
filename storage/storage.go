package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cgCodeLife/image_recognition_web/config"
	"github.com/cgCodeLife/image_recognition_web/errno"
	"github.com/cgCodeLife/image_recognition_web/videolib"
	"github.com/cgCodeLife/logs"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	DSN         string
	Active      int
	Idle        int
	IdleTimeout time.Duration
}

var (
	storage *Storage
)

type GCFrame struct {
	frameId   []string
	isRunning bool
	lock      sync.RWMutex
}

func NewGCFrame() *GCFrame {
	return &GCFrame{
		frameId:   make([]string, 0),
		isRunning: false,
	}
}
func (gc *GCFrame) Run() {
	gc.lock.RLock()
	if gc.isRunning {
		gc.lock.RUnlock()
		return
	}
	gc.lock.RUnlock()
	//set running
	gc.lock.Lock()
	gc.isRunning = true
	gc.lock.Unlock()
	for {
		frameId := make([]string, 0)
		for _, key := range gc.frameId {
			frame, ok := storage.Get(key)
			if !ok {
				frameId = append(frameId, key)
				continue
			}
			if expired := frame.Expired(); expired {
				logs.Infof("key=%s will delete", key)
				storage.Remove(key)
			}
		}
		gc.lock.Lock()
		gc.frameId = frameId
		gc.lock.Unlock()
		time.Sleep(60 * time.Second)
	}
}

func (gc *GCFrame) Store(key string) {
	gc.lock.Lock()
	defer gc.lock.Unlock()
	gc.frameId = append(gc.frameId, key)
}

var (
	gcFrame *GCFrame
)

//tag info
const (
	GET_FACE_VIDEO_INFO    = 1000
	GET_VIDEO_INFO         = 1001
	GET_VID                = 1002
	GET_VEHICLE_VIDEO_INFO = 2000
	GET_FACE_HISTORY       = 3000
	ADD_FACE_INFO          = 4000
	UPDATE_FACE_INFO       = 4001
	SEARCH_FACE_INFO       = 5000
	UPDATE_CAMERA_INFO     = 5001
	SEARCH_CAMERA_INFO     = 5002
	SEARCH_VEHICLE_INFO    = 5003
	SEARCH_CAMERA_UID      = 5004
)

const (
	FRAME_STEP = 1 //最大误差不超过1s,否则就认为没有合适的帧
)

const (
	VIDEO_PATH  = "_VIDEO_PATH"
	VIDEO_FRAME = "_FRAMES"
)

//table info
const (
	videoTable = "video_info"
	/**
	video_info | CREATE TABLE `video_info` (
	  `id` int(10) unsigned NOT NULL,
	  `vid` varchar(128) NOT NULL DEFAULT '000000000' COMMENT 'video id',
	  `uid` varchar(128) NOT NULL DEFAULT '111111111' COMMENT 'camera id',
	  `path` varchar(256) NOT NULL DEFAULT '/home/caoge/dev/src/github.com/cgCodeLife/image_recognition_web/output/data/face_video_1.mp4' COMMENT 'Video stores the local address',
	  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'video type 0:face 1:vehicle',
	  PRIMARY KEY (`id`),
	  KEY `uid` (`uid`),
	  CONSTRAINT `video_info_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `camera_info` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8
	*/
	cameraTable = "camera_info"
	/**
	 * table info:
	camera_info | CREATE TABLE `camera_info` (
	  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	  `uid` varchar(128) NOT NULL DEFAULT '111111111' COMMENT 'camera id',
	  `mac` varchar(64) NOT NULL DEFAULT '' COMMENT 'MAC address',
	  `addr` varchar(128) NOT NULL DEFAULT '128 zhichun road, zhongguancun, haidian district, Beijing, China',
	  `lat` float NOT NULL DEFAULT '39.9764' COMMENT 'latitude',
	  `lon` float NOT NULL DEFAULT '116.342' COMMENT 'longitude',
	  PRIMARY KEY (`id`),
	  UNIQUE KEY `uid_idx` (`uid`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 |
	*/
	faceLibTable = "face_info_lib"
	/**
	 CREATE TABLE `face_info_lib` (
	  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	  `pic` varchar(256) NOT NULL DEFAULT '' COMMENT '图像的base64编码信息',
	  `human_id` varchar(64) DEFAULT NULL,
	  `time` datetime NOT NULL,
	  `note` varchar(256) DEFAULT NULL,
	  `name` varchar(64) DEFAULT NULL,
	  PRIMARY KEY (`id`),
	  UNIQUE KEY `human_id` (`human_id`),
	  KEY `name` (`name`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	*/
	vehicleLibTable = "vehicle_info_lib"
	/**
	 CREATE TABLE `vehicle_info_lib` (
	  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	  `human_id` varchar(64) DEFAULT NULL,
	  `time` datetime NOT NULL,
	  `note` varchar(256) DEFAULT NULL,
	  `name` varchar(64) DEFAULT NULL,
	  `uid` varchar(64) DEFAULT NULL,
	  PRIMARY KEY (`id`),
	  UNIQUE KEY `uid` (`uid`),
	  KEY `human_id` (`human_id`),
	  CONSTRAINT `vehicle_info_lib_ibfk_1` FOREIGN KEY (`human_id`) REFERENCES `face_info_lib` (`human_id`) ON DELETE CASCADE ON UPDATE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	*/
	humanFaceRecordTable = "human_face_record_info"
	/* CREATE TABLE `human_face_record_info` (
	  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	  `vid` varchar(128) NOT NULL DEFAULT '' COMMENT 'video id',
	  `human_id` varchar(64) NOT NULL DEFAULT '000000000000000000',
	  `time` datetime NOT NULL,
	  `credibility` float NOT NULL DEFAULT '0',
	  `uid` varchar(128) NOT NULL DEFAULT '111111111' COMMENT 'camera id',
	  PRIMARY KEY (`id`),
	  KEY `uid` (`uid`),
	  KEY `human_id` (`human_id`),
	  CONSTRAINT `human_face_record_info_ibfk_2` FOREIGN KEY (`human_id`) REFERENCES `face_info_lib` (`human_id`) ON DELETE CASCADE ON UPDATE CASCADE,
	  CONSTRAINT `human_face_record_info_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `camera_info` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE
	) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8; */
)

type MemDB struct {
}

type Cache struct {
	sync.Map
}

type Storage struct {
	c     *Config
	db    *sql.DB
	mdb   *MemDB
	cache *Cache
}

func init() {
	if gcFrame != nil {
		return
	}
	gcFrame = NewGCFrame()
}
func (c *Cache) SearchFrame(vid, startTime, endTime string) (*videolib.Frame, error) {
	index, err := c.getIndex(startTime, endTime)
	if err != nil {
		return nil, err
	}
	key := c.GetFrameIndexKey(c.GetFramePrefixKey(vid), index)
	frame, ok := c.Load(key)
	if ok {
		return frame.(*videolib.Frame), nil
	}
	return nil, nil
}

func (c *Cache) getIndex(startTime, endTime string) (int, error) {
	//conver startTime and endTime to index
	begin, err := c.getTimeOffset(startTime)
	if err != nil {
		return -1, err
	}
	end, err := c.getTimeOffset(endTime)
	if err != nil {
		return -1, err
	}
	if end <= begin {
		return -1, errors.New("invalid time range")
	}
	index := c.rangeToIndex(begin, end)
	return index, nil
}

func (c *Cache) Search(key, startTime, endTime string) (*videolib.Frame, error) {
	//conver startTime and endTime to index
	index, err := c.getIndex(startTime, endTime)
	if err != nil {
		return nil, err
	}
	frame, err := c.find(key, index)
	if err != nil || frame == nil {
		return nil, err
	}
	return frame, nil
}

func (c *Cache) rangeToIndex(begin, end int) int {
	return (begin + end) / 2
}

func (c *Cache) PutFrame(vid, startTime, endTime string, frame *videolib.Frame) error {
	begin, err := c.getTimeOffset(startTime)
	if err != nil {
		logs.Errorf("method=PutFrame error=%s", err)
		return err
	}
	end, err := c.getTimeOffset(endTime)
	if err != nil {
		logs.Errorf("method=PutFrame error=%s", err)
		return err
	}
	if end <= begin {
		logs.Errorf("method=PutFrame error=%s", err)
		return errors.New("invalid time range")
	}
	index := c.rangeToIndex(begin, end)
	key := c.GetFrameIndexKey(c.GetFramePrefixKey(vid), index)
	c.Store(key, frame)
	gcFrame.Store(key)
	logs.Infof("frame key=%s start_time=%s end_time=%s success", key, startTime, endTime)
	return nil
}

func (c *Cache) GetFrameIndexKey(key string, index int) string {
	return key + "_" + strconv.Itoa(index)
}

func (s *Cache) GetFramePrefixKey(vid string) string {
	return VIDEO_FRAME + "_" + vid
}

func (c *Cache) getTimeOffset(t string) (int, error) {
	//t format: hh:mm:ss eg: 01:10:10
	tmp := strings.Split(t, ":")
	if len(tmp) != 3 {
		return -1, errors.New("invalid time")
	}
	h := tmp[0]
	m := tmp[1]
	s := tmp[2]
	digitHour, err := strconv.Atoi(h)
	if err != nil {
		return -1, err
	}
	hTos := 3600 * digitHour
	digitMinute, err := strconv.Atoi(m)
	if err != nil {
		return -1, err
	}
	mTos := 60 * digitMinute
	digitSecond, err := strconv.Atoi(s)
	if err != nil {
		return -1, err
	}
	return hTos + mTos + digitSecond, nil
}

func (c *Cache) find(key string, index int) (*videolib.Frame, error) {
	data, ok := c.Load(key + strconv.Itoa(index))
	if ok {
		return data.(*videolib.Frame), nil
	}
	//find frame before FRAME_STEP
	var res *videolib.Frame
	if index-FRAME_STEP >= 0 {
		if tmp, ok := c.Load(key + strconv.Itoa(index-FRAME_STEP)); ok {
			res = tmp.(*videolib.Frame)
			return res, nil
		}
	}
	//find frame in the future FRAME_STEP
	if tmp, ok := c.Load(key + strconv.Itoa(index+FRAME_STEP)); ok {
		res = tmp.(*videolib.Frame)
	}
	return res, nil
}

func DefaultMemDB() *MemDB {
	return &MemDB{}
}

func DefaultCache() *Cache {
	return &Cache{}
}

func new(c *Config) *Storage {
	return &Storage{
		c:     c,
		db:    newMySql(c),
		mdb:   DefaultMemDB(),
		cache: DefaultCache(),
	}
}

func newMySql(c *Config) *sql.DB {
	db, err := open(c)
	if err != nil {
		panic(err)
	}
	return db
}

func open(c *Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("sql.Open() error(%v)", err))
		return nil, err
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(c.IdleTimeout)
	return db, nil
}

func Init() {
	conf := &Config{
		DSN:         config.GetMysql().GetDSN(),
		Active:      config.GetMysql().GetActive(),
		Idle:        config.GetMysql().GetIdle(),
		IdleTimeout: config.GetMysql().GetIdleTimeout(),
	}
	storage = new(conf)
	//load vid to cache
	go storage.loadVids()
	//generate frames
	go storage.loadVideoFrames()
	go gcFrame.Run()
}

func (s *Storage) loadVids() {
	start := time.Now()
	vids := make([]string, 0)
	defer func() {
		logs.Infof("method=loadVids cost=%v load success, total=%d", time.Since(start), len(vids))
	}()
	stm, err := s.db.Prepare(fmt.Sprintf("select vid, path from %s", videoTable))
	if err != nil {
		logs.Errorf("method=loadVids load db error=%s", err)
		return
	}
	defer stm.Close()
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("method=loadVids db query error=%s", err)
		return
	}
	defer rows.Close()
	paths := make([]string, 0)
	for rows.Next() {
		var vid string
		var path string
		if err := rows.Scan(&vid, &path); err != nil {
			logs.Errorf("scan row error=%s", err)
		}
		vids = append(vids, vid)
		paths = append(paths, path)
	}
	//set cache
	for i := 0; i < len(vids); i++ {
		s.cache.Store(s.GetPathKey(vids[i]), paths[i])
	}
}

func (s *Storage) loadVideoFrames() {
	var total int
	start := time.Now()
	defer func() {
		logs.Infof("method=loadVideoFrames cost=%v load success, total=%d", time.Since(start), total)
	}()
	stm, err := s.db.Prepare(fmt.Sprintf("select vid, path from %s", videoTable))
	if err != nil {
		logs.Errorf("method=loadVideoFrames load db error=%s", err)
		return
	}
	defer stm.Close()
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("method=loadVideoFrames db query error=%s", err)
		return
	}
	defer rows.Close()
	vids := make([]string, 0)
	paths := make([]string, 0)
	for rows.Next() {
		var vid string
		var path string
		if err := rows.Scan(&vid, &path); err != nil {
			logs.Errorf("scan row error=%s", err)
		}
		vids = append(vids, vid)
		paths = append(paths, path)
	}
	//set cache
	for i := 0; i < len(vids); i++ {
		//you can concurrent running
		s.cacheAllFrame(vids[i], paths[i])
		total += 1
	}
}

func (s *Storage) cacheAllFrame(vid, path string) {
	begin := "00:00:00"
	end := begin
	for {
		end = s.addTime(begin, FRAME_STEP)
		logs.Infof("begin=%s end=%s", begin, end)
		frame, err := videolib.GetFrame(path, begin, end)
		if err != nil {
			logs.Errorf("path=%s begin=%s end=%s vid=%s get frame error=%s", path, begin, end, vid, err)
			continue
		}
		if frame == nil {
			break
		}
		s.cache.PutFrame(vid, begin, end, frame)
		begin = end
	}
}

func (s *Storage) close() {
	s.db.Close()
}

func (s *Storage) addTime(begin string, incr int) string {
	//eg: 00:00:00 + 3 ==> 00:00:03
	tmp := strings.Split(begin, ":")
	if len(tmp) != 3 {
		logs.Errorf("begin=%s format error you should input hh:mm:ss", begin)
		return "00:00:03"
	}
	hour := tmp[0]
	minute := tmp[1]
	second := tmp[2]
	digitHour, err := strconv.Atoi(hour)
	if err != nil {
		logs.Errorf("begin=%s get hour error=%s", begin, err)
		return begin
	}

	digitMinute, err := strconv.Atoi(minute)
	if err != nil {
		logs.Errorf("begin=%s get minute error=%s", begin, err)
		return begin
	}
	digitSecond, err := strconv.Atoi(second)
	if err != nil {
		logs.Errorf("begin=%s get second error=%s", begin, err)
		return begin
	}
	digitSecond += incr
	cap := digitSecond / 60
	digitSecond = digitSecond % 60
	if cap != 0 {
		digitMinute += cap
		cap = digitMinute / 60
		digitMinute = digitMinute % 60

		if cap != 0 {
			digitHour += cap
		}
	}
	if digitSecond < 10 {
		second = fmt.Sprintf("0%d", digitSecond)
	} else {
		second = fmt.Sprintf("%d", digitSecond)
	}
	if digitMinute < 10 {
		minute = fmt.Sprintf("0%d", digitMinute)
	} else {
		minute = fmt.Sprintf("%d", digitMinute)
	}
	if digitHour < 10 {
		hour = fmt.Sprintf("0%d", digitHour)
	} else {
		hour = fmt.Sprintf("%d", digitHour)
	}
	return fmt.Sprintf("%s:%s:%s", hour, minute, second)
}

func (s *Storage) dbQueryWrite(ctx context.Context, sql string) error {
	start := time.Now()
	defer func() {
		logs.CtxInfo(ctx, "sql=%s cost=%s", sql, time.Since(start))
	}()
	tx, _ := s.db.Begin()
	defer tx.Rollback()

	res, err := tx.Exec(sql)
	if err != nil {
		logs.CtxError(ctx, "sql=%s exec error=%s", err)
		return err
	}
	logs.CtxInfo(ctx, "sql=%s exec result=%v", res)
	affecRows, err := res.RowsAffected()
	if err != nil {
		logs.CtxError(ctx, "RowsAffected error=%s", err)
		return err
	}
	if affecRows == int64(0) {
		logs.CtxWarn(ctx, "sql=%s not found", sql)
		return errno.NOT_FOUND
	}
	err = tx.Commit()
	return err
}

func (s *Storage) dbQuery(ctx context.Context, tag int, sql string) (interface{}, error) {
	start := time.Now()
	defer func() {
		logs.CtxInfo(ctx, "sql=%s cost=%v", sql, time.Since(start))
	}()
	stm, err := s.db.Prepare(sql)
	if err != nil {
		logs.CtxError(ctx, "Prepare error=%s", err)
		return nil, err
	}
	defer stm.Close()
	switch tag {
	case GET_FACE_VIDEO_INFO:
		return getFaceVideoInfo(stm)
	case GET_VEHICLE_VIDEO_INFO:
		return getVehicleVideoInfo(stm)
	case GET_VIDEO_INFO:
		return getVideoInfo(stm)
	case GET_VID:
		return getVid(stm)
	case GET_FACE_HISTORY:
		return getFaceHistory(stm)
	case SEARCH_FACE_INFO:
		return searchFaceInfo(stm)
	case SEARCH_CAMERA_INFO:
		return searchCameraInfo(stm)
	case SEARCH_VEHICLE_INFO:
		return searchVehicleInfo(stm)
	case SEARCH_CAMERA_UID:
		return searchCameraUid(stm)
	default:
		logs.CtxError(ctx, "not found method for sql=%s", sql)
	}
	return nil, errors.New("not found")
}

func (s *Storage) SearchFrame(vid, startTime, endTime string) *videolib.Frame {
	data, err := s.cache.SearchFrame(vid, startTime, endTime)
	if err != nil {
		return nil
	}
	return data
}

func (s *Storage) Get(key string) (*videolib.Frame, bool) {
	frame, ok := s.cache.Load(key)
	if !ok {
		return nil, false
	}
	return frame.(*videolib.Frame), true
}

func (s *Storage) Remove(key string) {
	s.cache.Delete(key)
}

func (s *Storage) CacheFrame(vid, startTime, endTime string, frame *videolib.Frame) {
	s.cache.PutFrame(vid, startTime, endTime, frame)
}

func (s *Storage) GetPathKey(vid string) string {
	return VIDEO_PATH + "_" + vid
}

func (s *Storage) GetVideoFile(vid string) string {
	if res, ok := s.cache.Load(s.GetPathKey(vid)); ok {
		return res.(string)
	}
	return ""
}

func (s *Storage) GetUpdateVideoFile(vid string) (string, error) {
	//load from db
	sql := fmt.Sprintf("select vid from %s where vid=\"%s\"", videoTable, vid)
	file, err := s.dbQuery(context.TODO(), GET_VID, sql)
	if err != nil {
		return "", err
	}
	//update cache
	if file != nil {
		s.cache.Store(s.GetPathKey(vid), file)
	}
	return file.(string), nil
}

func getVid(stm *sql.Stmt) (string, error) {
	file := ""
	rows, err := stm.Query()
	if err != nil {
		logs.Errorf("query error=%s", err)
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&file); err != nil {
			logs.Errorf("scan row error=%s", err)
			return "", err
		}
	}
	return file, nil
}
