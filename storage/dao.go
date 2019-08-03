package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cgCodeLife/image_recognition_web/config"
	"github.com/cgCodeLife/logs"
	_ "github.com/go-sql-driver/mysql"
	"os"
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

//tag info
const (
	GET_FACE_VIDEO_INFO = 1000
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
	humanFaceRecordTable = "human_face_record_info"
	/**
	 * table info:
	 human_face_record_info | CREATE TABLE `human_face_record_info` (
	  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	  `vid` varchar(128) NOT NULL DEFAULT '' COMMENT 'video id',
	  `human_id` varchar(64) NOT NULL DEFAULT '000000000000000000',
	  `time` datetime NOT NULL,
	  `credibility` float NOT NULL DEFAULT '0',
	  PRIMARY KEY (`id`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 |
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
)

type MemDB struct {
}

type Cache struct {
}

type Storage struct {
	c     *Config
	db    *sql.DB
	mdb   *MemDB
	cache *Cache
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
}

func (s *Storage) close() {
	s.db.Close()
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
	default:
		logs.CtxError(ctx, "not found method for sql=%s", sql)
	}
	return nil, errors.New("not found")
}
