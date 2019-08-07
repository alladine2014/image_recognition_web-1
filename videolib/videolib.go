package videolib

import (
	// "errors"
	"github.com/cgCodeLife/logs"
	"io/ioutil"
	"os"
	"os/exec"
)

type Frame []byte

func GetFrame(file, startTime, endTime string) ([]Frame, error) {
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
	return []Frame{Frame(data)}, nil
}
