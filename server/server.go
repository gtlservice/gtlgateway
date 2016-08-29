package server

import "github.com/gtlservice/gtlgateway/etc"
import "github.com/gtlservice/gtlgateway/ctrl"
import "github.com/gtlservice/gtlgateway/cache"
import "github.com/gtlservice/gutils/logger"
import "github.com/gtlservice/gutils/system"
import "github.com/gtlservice/gzkwrapper"

import (
	"flag"
	"net/http"
)

type Gateway struct {
	GatewayHandler
	Configuration *etc.Configuration
	RouteCache    *cache.RouteCache
	Controller    *ctrl.Controller
}

func NewGateway() (*Gateway, error) {

	path, err := system.GetExecDir()
	if err != nil {
		return nil, err
	}

	key, err := system.MakeKeyFile(path + "/gtlgateway.key")
	if err != nil {
		return nil, err
	}

	var etcfile string
	flag.StringVar(&etcfile, "f", "etc/config.yaml", "gateway etc file.")
	flag.Parse()
	configuration, err := etc.NewConfiguration(etcfile)
	if err != nil {
		return nil, err
	}

	gateway := &Gateway{}
	largs := configuration.GetLogger()
	logger.OPEN(largs)
	zkargs := configuration.GetZkWrapper()
	server, err := gzkwrapper.NewServer(key, zkargs, gateway)
	if err != nil {
		return nil, err
	}

	route := cache.NewRouteCache(gateway)
	controller := ctrl.NewController(server, route)
	gateway.Configuration = configuration
	gateway.RouteCache = route
	gateway.Controller = controller
	return gateway, nil
}

func (gateway *Gateway) Startup() error {

	if err := gateway.Controller.Initialize(); err != nil {
		logger.ERROR("[#server#] gateway controller initialize error, %s", err)
		return err
	}

	go func() {
		bind := gateway.Configuration.GetHttpBind()
		if err := http.ListenAndServe(bind, gateway); err != nil {
			logger.ERROR("[#server#] gateway start error:%s", err.Error())
		}
	}()

	logger.INFO("[#server#] gateway started.")
	logger.INFO("[#server#] gateway configuration...\r\n%v", gateway.Configuration)
	return nil
}

func (gateway *Gateway) Stop() error {

	if err := gateway.Controller.UnInitialize(); err != nil {
		logger.ERROR("[#server#] gateway controller uninitialize error, %s", err)
		return err
	}

	logger.INFO("[#server#] gateway closed.")
	logger.CLOSE()
	return nil
}
