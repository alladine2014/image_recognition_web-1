package errno

import (
	"errors"
	"fmt"
)

// error codes
// 30xxx，网关分配的前两位，成功统一返回2000
var (
	// 4xx
	InvalidParam        = Payload{Code: 30402, Message: "invalid json param"}
	OutOfBounds         = Payload{Code: 30405, Message: "object_ids exceeded"}
	InvalidToken        = Payload{Code: 30406, Message: "invalid token"}
	ProcessTimeout      = Payload{Code: 30410, Message: "image process time out"}
	TooManyRequest      = Payload{Code: 30411, Message: "too many reqeust"}
	InvalidTopAccount   = Payload{Code: 30414, Message: "invalid top Account"}
	InvalidTestHost     = Payload{Code: 30415, Message: "invalid test host"}
	InvalidVid          = Payload{Code: 30416, Message: "invalid vid"}
	InvalidTime         = Payload{Code: 30417, Message: "invalid start_time or end_time"}
	InvalidFaceId       = Payload{Code: 30418, Message: "invalid face id"}
	InvalidVehicleId    = Payload{Code: 30419, Message: "invalid vehicle id"}
	InvalidPic          = Payload{Code: 30420, Message: "invalid picture"}
	InvalidHumanName    = Payload{Code: 30421, Message: "invalid name"}
	InvalidDate         = Payload{Code: 30422, Message: "invalid time you should input this format: yyyy-MM-dd HH:mm:ss"}
	InvalidFaceIdOrName = Payload{Code: 30423, Message: "invalid name or id"}
	InvalidCameraId     = Payload{Code: 30424, Message: "invalid camera id"}
	InvalidMac          = Payload{Code: 30425, Message: "invalid mac"}
	InvalidAddr         = Payload{Code: 30426, Message: "invalid addr"}
	InvalidLocation     = Payload{Code: 30427, Message: "invalid location"}
	ParseJsonError      = Payload{Code: 30428, Message: "parse json error"}
	NotFoundRecord      = Payload{Code: 30429, Message: "not found record"}
	// 5xx
	InternalErr    = Payload{Code: 30500, Message: "internal error"}
	GetFileErr     = Payload{Code: 30502, Message: "get file error"}
	EncryptionErr  = Payload{Code: 30509, Message: "entryption process error"}
	ActionNotFound = Payload{Code: 30510, Message: "action not found"}
	RunFunctionErr = Payload{Code: 30511, Message: "run function error"}
	SaveDBErr      = Payload{Code: 30512, Message: "save db error"}
	SearchDBErr    = Payload{Code: 30513, Message: "query db error"}
	AddDBErr       = Payload{Code: 30514, Message: "insert db error"}
	UpdateDBErr    = Payload{Code: 30514, Message: "update db error"}
	DelDBErr       = Payload{Code: 30514, Message: "delete db error"}
)

var (
	INVALID_VID = errors.New("InvalidVid")
	NOT_FOUND   = errors.New("NotFound")
)

// Payload defines http body for response
type Payload struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	File    string      `json:"-"`
}

// OK response success case
func OK(data interface{}) Payload {
	return Payload{
		Code:    2000,
		Message: "success",
		Data:    data,
	}
}

func LocalStream(name string) Payload {
	return Payload{Code: 2000, Message: "success", File: name}
}

// Error make paylaod support error type
func (p Payload) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", p.Code, p.Message)
}

func InternalError(err error) Payload {
	p := InternalErr
	p.Message = err.Error()
	return p
}
