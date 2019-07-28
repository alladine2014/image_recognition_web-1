package main

import (
	"fmt"
	"github.com/cgCodeLife/image_recognition_web/config"
	"github.com/cgCodeLife/image_recognition_web/storage"
	"github.com/gin-gonic/gin"
	"os"
	"runtime"
)

func init() {
	os.Setenv("GODEBUG", fmt.Sprintf("netdns=cgo,%s", os.Getenv("GODEBUG")))
	config.LoadConf()
	gin.SetMode("debug")
	runtime.GOMAXPROCS(runtime.NumCPU())
	storage.Init()
}

func main() {
}
