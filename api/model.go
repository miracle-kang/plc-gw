package api

import (
	"time"
)

// Request command model
type ReadTagCommand struct {
	SN         string   `json:"SN"`
	DeviceName string   `json:"DeviceName"`
	Data       []string `json:"Data"`
}

type WriteTagCommand struct {
	SN         string      `json:"SN"`
	DeviceName string      `json:"DeviceName"`
	WTagName   string      `json:"WTagName"`
	WTagValue  interface{} `json:"WTagValue"`
}

type ClearGatewayCommand struct {
	SN          string   `json:"SN"`
	DeviceNames []string `json:"DeviceNames"`
}

type ReportDataCommand struct {
	Time         string                 `json:"Time"`
	SN           string                 `json:"SN"`
	DeviceName   string                 `json:"DeviceName"`
	GetDataError interface{}            `json:"GetDataError"`
	Data         map[string]interface{} `json:"Data"`
}

// Response model
type BaseResponse struct {
	Result  bool        `json:"Result"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

type ReportDataResponse struct {
	DeviceName string                 `json:"DeviceName"`
	Data       map[string]interface{} `json:"Data"`
}

type Empty struct{}

type TagDto struct {
	Name  string      `json:"Name"`
	Value interface{} `json:"Value"`
}

type GatewayDto struct {
	SN     string   `json:"SN"`
	IP     string   `json:"IP"`
	Online bool     `json:"Online"`
	PLCs   []string `json:"PLCs"`
}

type PLCDto struct {
	SN         string     `json:"SN"`
	Name       string     `json:"Name"`
	ConnState  bool       `json:"ConnState"`
	RTags      []*TagDto  `json:"RTags"`
	WTags      []*TagDto  `json:"WTags"`
	LastReport *time.Time `json:"LastReport"`
	LastWrite  *time.Time `json:"LastWrite"`
}

func Failed(message string) BaseResponse {
	return FailedWithData(message, &Empty{})
}

func FailedWithData(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Result:  false,
		Message: message,
		Data:    data,
	}
}

func Success(message string) BaseResponse {
	return SuccessWithData(message, &Empty{})
}

func SuccessWithData(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Result:  true,
		Message: message,
		Data:    data,
	}
}
