package videolib

import (
	// "errors"
	"github.com/cgCodeLife/logs"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type Frame struct {
	Data       []byte
	TTL        int
	CreateTime int64
}

func NewFrame(data []byte) *Frame {
	return &Frame{
		Data:       data,
		TTL:        600,
		CreateTime: time.Now().Unix(),
	}
}

func (fr *Frame) Expired() bool {
	if time.Now().Unix()-fr.CreateTime >= int64(fr.TTL) {
		logs.Infof("test time ttl: diff=%d ttl=%d", time.Now().Unix()-fr.CreateTime, fr.TTL)
		return true
	}
	return false
}

func (fr *Frame) GetData() []byte {
	return fr.Data
}

func GetFrame(file, startTime, endTime string) (*Frame, error) {
	os.Remove("image.png")
	//version 1.0 use cmd directly
	cmd := exec.Command("ffmpeg", "-i", file, "-ss", startTime, "-t", endTime, "-r", "1", "image.png")
	if err := cmd.Run(); err != nil {
		logs.Errorf("cmd result error=%s", err)
		// return nil, errors.New("cloud not generate frame")
	}
	data, err := ioutil.ReadFile("image.png")
	if err != nil {
		logs.Errorf("read frame file image.png error=%s", err)
		return nil, nil
	}
	logs.Infof("read data size=%d", len(data))
	return NewFrame(data), nil
}
