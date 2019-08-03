package middleware

import (
	"context"
	"encoding/hex"
	"github.com/cgCodeLife/image_recognition_web/errno"
	"github.com/cgCodeLife/image_recognition_web/util"
	"github.com/cgCodeLife/logs"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"
)

const (
	TAG_TRACEID = "X-Tt-Traceid"
	TAG_PSM     = "X-Psm"
	TAG_LOGID   = "K_LOGID"
)

type MyHandler func(c *gin.Context, ctx context.Context) errno.Payload

func Access() gin.HandlerFunc {
	return func(c *gin.Context) {
		mkey := "access"
		tagkv := map[string]string{"path": c.Request.URL.Path}
		util.EmitThroughput(mkey, tagkv)
		defer util.EmitLatency(createContext(getLogID(c)), mkey, time.Now(), tagkv)

		logid := getLogID(c)
		c.Writer.Header().Set("X-Ai-Requestid", logid)

		c.Next()
	}
}

func Response(mkey string, f MyHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := createContext(getLogID(c))
		tagkv := make(map[string]string)
		defer util.EmitLatency(ctx, mkey, time.Now(), tagkv)
		defer util.EmitThroughput(mkey, tagkv)
		ctx = logs.CtxAddKVs(ctx, "method", c.Request.Method)
		ctx = logs.CtxAddKVs(ctx, "url", c.Request.URL.String())
		ctx = logs.CtxAddKVs(ctx, "remote_addr", c.Request.RemoteAddr)
		data := f(c, ctx)

		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-TT-Access, Content-Type, accept, content-disposition, content-range")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// disable any cache, https://stackoverflow.com/questions/49547/how-to-control-web-page-caching-across-all-browsers
		c.Writer.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate, proxy-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache") // http 1.0
		c.Writer.Header().Set("Expires", "0")
		c.JSON(http.StatusOK, data)
		logs.CtxInfo(ctx, "request=%+v, response=%+v", c.Request, data)
		tagkv["code"] = strconv.Itoa(data.Code)
		if data.Code != errno.OK(nil).Code {
			util.EmitError(mkey, tagkv)
		}
	}
}

func createContext(logID string) context.Context {
	ctx := context.Background()
	ctx = newCtxWithLogID(ctx, logID)
	return ctx
}

func getLogID(c *gin.Context) string {
	logID := c.Request.Header.Get("X-TT-LogID")
	if c.Request.Header.Get("X-TT-TraceID") != "" {
		logID = c.Request.Header.Get("X-TT-TraceID")
	}
	if c.Request.Header.Get("X-Top-Request-Id") != "" {
		logID = c.Request.Header.Get("X-Top-Request-Id")
	}
	if len(logID) == 0 {
		logID = getID()
	}
	return logID
}

func newCtxWithLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, TAG_LOGID, logID)
}

func getID() string {
	temp, err := uuid.NewV4()
	if err != nil {
		logs.Errorf("uuid generate error=%s", err)
		return ""
	}
	return hex.EncodeToString(temp.Bytes())
}
