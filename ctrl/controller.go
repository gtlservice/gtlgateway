package ctrl

import "github.com/gtlservice/gtlgateway/cache"
import "github.com/gtlservice/gutils/logger"
import "github.com/gtlservice/gzkwrapper"

type Controller struct {
	Server *gzkwrapper.Server
	Route  *cache.RouteCache
}

func NewController(server *gzkwrapper.Server, route *cache.RouteCache) *Controller {

	return &Controller{
		Server: server,
		Route:  route,
	}
}

func (c *Controller) Initialize() error {

	logger.INFO("[#ctrl#] controller initializeing......")
	if err := c.Server.Open(); err != nil {
		return err
	}
	return nil
}

func (c *Controller) UnInitialize() error {

	c.Route.ClearServer()
	if err := c.Server.Close(); err != nil {
		return err
	}
	return nil
}
