package base

import (
	"errors"
)

type Service struct {
	Name string `json:"name,omitempty"` //服务名称
	Host string `json:"host,omitempty"` //服务地址
}

func ReadServiceData(v interface{}) (*Service, error) {

	switch p := v.(type) {
	case map[string]interface{}:
		return &Service{
			Name: p["name"].(string),
			Host: p["host"].(string),
		}, nil
	}
	return nil, errors.New("service data invalid.")
}

type Server struct {
	Key         string `json:"key"`         //服务器唯一编码
	ServiceName string `json:"servicename"` //服务名称
	HostName    string `json:"hostname"`    //主机名称
	Location    string `json:"location"`    //节点位置
	OS          string `json:"os"`          //系统信息
	Platform    string `json:"platform"`    //平台信息
	Host        string `json:"host"`        //服务地址
}

func NewServer(key string, servicename string, hostname string, location string,
	os string, platform string, host string) *Server {

	return &Server{
		Key:         key,
		ServiceName: servicename,
		HostName:    hostname,
		Location:    location,
		OS:          os,
		Platform:    platform,
		Host:        host,
	}
}
