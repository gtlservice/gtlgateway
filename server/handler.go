package server

import "github.com/gtlservice/gtlgateway/base"
import "github.com/gtlservice/gtlgateway/cache"
import "github.com/gtlservice/gutils/logger"
import "github.com/gtlservice/gutils/system"
import "github.com/gtlservice/gzkwrapper"

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type GatewayHandler interface {
	http.Handler
	cache.IRouteCacheHandler
	gzkwrapper.INodeNotifyHandler
}

func (gateway *Gateway) OnZkWrapperNodeHandlerFunc(online []*gzkwrapper.NodeInfo, offline []*gzkwrapper.NodeInfo) {

	logger.INFO("[#server#] zkwrapper online nodes:%d", len(online))
	for _, node := range online {
		service, err := base.ReadServiceData(node.Data.Attach)
		if err != nil {
			logger.ERROR("[#server#] node info error, %s:%s %s", node.Key, node.Data.IpAddr, err)
		} else {
			logger.INFO("[#server#] node info %s:%s %s:%s", node.Key, node.Data.IpAddr, service.Name, service.Host)
			gateway.RouteCache.CreateServer(service.Name, node.Key, node.Data.HostName, node.Data.Location, node.Data.OS, node.Data.Platform, service.Host)
		}
	}

	logger.INFO("[#server#] zkwrapper offline nodes:%d", len(offline))
	for _, node := range offline {
		service, err := base.ReadServiceData(node.Data.Attach)
		if err != nil {
			logger.ERROR("[#server#] node info error, %s:%s %s", node.Key, node.Data.IpAddr, err)
		} else {
			logger.INFO("[#server#] node info %s:%s %s:%s", node.Key, node.Data.IpAddr, service.Name, service.Host)
			gateway.RouteCache.RemoveServer(service.Name, node.Key)
		}
	}
}

func (gateway *Gateway) OnZkWrapperPulseHandlerFunc(key string, nodedata *gzkwrapper.NodeData, err error) {
	//可不实现该回调
}

func (gateway *Gateway) OnZkWrapperWatchHandlerFunc(path string, data []byte, err error) {
	//可不实现该回调
}

func (gateway *Gateway) OnRouteCacheChangedHandleFunc(event cache.EventType, server *base.Server) {

	logger.INFO("[#server#] route cache change %s:%s %s", event.String(), server.Key, server.ServiceName)
}

func (gateway *Gateway) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	}

	if req.Method == "OPTIONS" {
		return
	}

	uri := strings.TrimSpace(req.RequestURI)
	if uri == "/" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid gateway request."))
		return
	}

	array := strings.Split(uri, "/")
	server := gateway.RouteCache.HashServer(array[1], system.MakeKey(true))
	if server == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("route failed, request cannot be transfer."))
		return
	}

	uri = fmt.Sprintf("http://%s%s", server.Host, uri)
	logger.INFO("[#server#] servehttp transfer:%s:%s", req.Method, uri)
	code, result, err := transToHttp(uri, req)
	if err != nil {
		result = []byte(err.Error())
	}
	w.WriteHeader(code)
	w.Write(result)
}

func transToHttp(uri string, req *http.Request) (int, []byte, error) {

	u, err := url.Parse(uri)
	if err != nil {
		return http.StatusBadRequest, []byte{}, err
	}

	temp := *req
	new_req := &temp
	new_req.URL = u
	new_req.Host = u.Host
	new_req.RequestURI = ""
	client := &http.Client{}
	resp, err := client.Do(new_req)
	if err != nil {
		return http.StatusBadRequest, []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusBadRequest, []byte{}, err
	}
	fmt.Println(resp.StatusCode)
	return resp.StatusCode, body, nil
}
