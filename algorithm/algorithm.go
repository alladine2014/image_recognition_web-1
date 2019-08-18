package algorithm

import (
	"bytes"
	"encoding/json"
	"github.com/cgCodeLife/image_recognition_web/config"
	"github.com/cgCodeLife/image_recognition_web/videolib"
	"io/ioutil"
	"net/http"
)

var (
	client *http.Client
	host   string
)

func Init() {
	client = &http.Client{}
	host = config.GetAlgorithmHost()
}

type FrameFaceRes struct {
	Body    BodyInfo `json:"body"`
	Message string   `json:"message"`
	Status  string   `json:"status"`
}

type BodyInfo struct {
	Face    []FaceSt `face`
	imageId string   `image_id`
}

type FaceSt struct {
	BoudingBox     []float64 `json:"bounding_box"`
	Classification string    `classification`
}

func GetFrameVehicleInfo(frames []videolib.Frame) (string, error) {
	return "test data", nil
}

func GetFrameFaceInfo(frames []videolib.Frame) (FrameFaceRes, error) {
	//debug
	return FrameFaceRes{
		Body: BodyInfo{
			Face: []FaceSt{
				{
					BoudingBox: []float64{
						1.1,
						2.2,
						3.3,
						4.4,
					},
					Classification: "harden",
				},
			},
		},
	}, nil
	res := FrameFaceRes{}
	url := host + "/pic_feed"
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(frames[0]))) //just handle first pic
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(body, &res); err != nil {
		return res, err
	}
	return res, nil
}
