package main

import (
	"fmt"
	"github.com/cgCodeLife/image_recognition_web/algorithm"
	"github.com/cgCodeLife/image_recognition_web/config"
	"github.com/cgCodeLife/image_recognition_web/ginex"
	"github.com/cgCodeLife/image_recognition_web/handle"
	"github.com/cgCodeLife/image_recognition_web/logger"
	"github.com/cgCodeLife/image_recognition_web/middleware"
	"github.com/cgCodeLife/image_recognition_web/storage"
	"github.com/cgCodeLife/logs"
	"os"
	"runtime"
)

func init() {
	os.Setenv("GODEBUG", fmt.Sprintf("netdns=cgo,%s", os.Getenv("GODEBUG")))
	config.LoadConf()
}

func main() {
	logger.Init()
	defer logger.Stop()
	ginex.SetMode("debug")
	runtime.GOMAXPROCS(runtime.NumCPU())
	storage.Init()
	algorithm.Init()

	r := ginex.Default()
	r.OPTIONS("/image_recognition/v1/*path", handle.Options)
	r.Use(middleware.Access())
	g := r.Group("/image_recognition/v1")
	g.GET("/face_video_info", middleware.Response("face.video.info", handle.GetFaceVideoInfo))
	g.GET("/video_info", middleware.Response("video.info", handle.GetVideoInfo))
	g.GET("/frame/face_info", middleware.Response("frame.face.info", handle.GetFrameFaceInfo))
	g.GET("/vehicle/video_info", middleware.Response("vehicle.video.info", handle.GetVehicleVideoInfo))
	g.GET("/frame/vehicle/info", middleware.Response("frame.vehicle.info", handle.GetFrameVehicleInfo))
	g.GET("/frame/vehicle/tracffic_flow", middleware.Response("frame.vehicle.traffic.flow",
		handle.GetFrameVehicleTrafficFlow))
	g.GET("/frame/vehicle/avg_speed", middleware.Response("frame.vehicle.avg.speed",
		handle.GetFrameVehicleAvgSpeed))
	g.GET("/face_history", middleware.Response("face.history", handle.GetFaceHistory))
	g.POST("/face_info", middleware.Response("add.face.info", handle.AddFaceInfo))
	g.GET("/face_info", middleware.Response("search.face.info", handle.SearchFaceInfo))
	g.PUT("/face_info", middleware.Response("update.face.info", handle.UpdateFaceInfo))
	g.DELETE("/face_info", middleware.Response("delete.face.info", handle.DelFaceInfo))
	g.POST("/camera_info", middleware.Response("add.camera.info", handle.AddCameraInfo))
	g.PUT("/camera_info", middleware.Response("update.camera.info", handle.UpdateCameraInfo))
	g.GET("/camera_info", middleware.Response("search.camera.info", handle.SearchCameraInfo))
	g.POST("/vehicle_info", middleware.Response("add.vehicle.info", handle.AddVehicleInfo))
	g.GET("/vehicle_info", middleware.Response("search.vehicle.info", handle.SearchVehicleInfo))
	g.PUT("/vehicle_info", middleware.Response("search.vehicle.info", handle.UpdateVechicleInfo))

	if err := r.Run(fmt.Sprintf("%s:%d", config.GetAddr(), config.GetPort())); err != nil {
		logs.Fatalf("server exit with error=%s", err)
		return
	}
	logs.Infof("server closed")
}
