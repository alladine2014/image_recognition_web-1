package util

import (
	"time"
)

func GetTimeNow() string {
	//return yyy-mm-dd HH:MM:SS
	return time.Now().Format("2006-01-02 15:04:05")
}
