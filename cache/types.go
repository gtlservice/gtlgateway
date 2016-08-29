package cache

import "github.com/gtlservice/gtlgateway/base"

type EventType int

const (
	CREATE_SERVER  EventType = iota + 1 //构建*base.Server事件
	RELEASE_SERVER                      //释放*base.Server事件
)

func (t EventType) String() string {
	switch t {
	case CREATE_SERVER:
		return "CREATE_SERVER"
	case RELEASE_SERVER:
		return "RELEASE_SERVER"
	}
	return ""
}

type IRouteCacheHandler interface {
	OnRouteCacheChangedHandleFunc(event EventType, server *base.Server)
}

type RouteCacheChangedHandleFunc func(event EventType, server *base.Server)

func (fn RouteCacheChangedHandleFunc) OnRouteCacheChangedHandleFunc(event EventType, server *base.Server) {
	fn(event, server)
}
