package algorithm

import (
	"bytes"
	"encoding/json"
	"github.com/cgCodeLife/image_recognition_web/config"
	"github.com/cgCodeLife/image_recognition_web/videolib"
	"github.com/cgCodeLife/logs"
	"io"
	"io/ioutil"
	"mime/multipart"
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

func GetFrameVehicleInfo(frame *videolib.Frame) (string, error) {
	return "test data", nil
}

func GetFrameFaceInfo(frame *videolib.Frame) (FrameFaceRes, error) {
	res := FrameFaceRes{}
	requrl := "http://" + host + "/pic_feed"
	binary := []byte(frame.GetData())
	//debug
	ioutil.WriteFile("test.png", binary, 0644)
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("file", "face.png")
	if err != nil {
		logs.Errorf("CreateFormFile error=%s", err)
		return res, err
	}
	_, err = io.Copy(formFile, bytes.NewReader(binary))
	if err != nil {
		logs.Errorf("copy frame to formFile error=%s", err)
		return res, err
	}
	contentType := writer.FormDataContentType()
	writer.Close()
	if err := writer.WriteField("id", "heyibing"); err != nil {
		logs.Errorf("WriteField id=heyibing error=%s", err)
		return res, err
	}
	resp, err := http.Post(requrl, contentType, buf)
	if err != nil {
		return res, err
	}
	// req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(frames[0]))) //just handle first pic
	// if err != nil {
	// 	return res, err
	// }
	// req.Header.Set("Content-Type", "multipart/form-data")
	// resp, err := client.Do(req)
	// if err != nil {
	// 	logs.Errorf("send request error=%s", err)
	// 	return res, err
	// }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Errorf("read response error=%s", err)
		return res, err
	}
	logs.Infof("response debug info:%s", string(body))
	if err = json.Unmarshal(body, &res); err != nil {
		logs.Errorf("Unmarshal error=%s", err)
		return res, err
	}
	return res, nil
}
