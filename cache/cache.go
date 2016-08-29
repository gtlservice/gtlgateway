package cache

import "github.com/gtlservice/gtlgateway/base"
import "github.com/gtlservice/gutils/algorithm"

import (
	"strings"
	"sync"
)

type RouteCache struct {
	sync.RWMutex
	servers map[string][]*base.Server
	circle  map[string]*algorithm.Consistent
	handler IRouteCacheHandler
}

func NewRouteCache(handler IRouteCacheHandler) *RouteCache {

	return &RouteCache{
		servers: make(map[string][]*base.Server, 0),
		circle:  make(map[string]*algorithm.Consistent, 0),
		handler: handler,
	}
}

func (cache *RouteCache) GetServer(servicename string, key string) *base.Server {

	cache.RLock()
	defer cache.RUnlock()
	if strings.TrimSpace(servicename) != "" {
		if _, ret := cache.servers[servicename]; ret {
			for _, server := range cache.servers[servicename] {
				if server.Key == key {
					return server
				}
			}
		}
	} else {
		for _, servers := range cache.servers {
			for _, server := range servers {
				if server.Key == key {
					return server
				}
			}
		}
	}
	return nil
}

func (cache *RouteCache) HashServer(servicename string, hkey string) *base.Server {

	cache.RLock()
	defer cache.RUnlock()
	if _, ret := cache.servers[servicename]; ret {
		ret := cache.circle[servicename].Get(hkey)
		servers := cache.servers[servicename]
		for _, server := range servers {
			if server.Key == ret {
				return server
			}
		}
	}
	return nil
}

func (cache *RouteCache) CreateServer(servicename string, key string, hostname string, location string, os string, platform string, host string) int {

	cache.Lock()
	defer cache.Unlock()
	server := base.NewServer(key, servicename, hostname, location, os, platform, host)
	if _, ret := cache.servers[servicename]; !ret {
		cache.servers[servicename] = []*base.Server{server}
		consistent := algorithm.NewConsisten(0xff)
		consistent.Add(server.Key)
		cache.circle[servicename] = consistent
	} else {
		found := false
		for _, value := range cache.servers[servicename] {
			if value.Key == server.Key {
				value = server
				found = true
				break
			}
		}
		if !found {
			cache.servers[servicename] = append(cache.servers[servicename], server)
			cache.circle[servicename].Add(server.Key)
		}
		cache.handler.OnRouteCacheChangedHandleFunc(CREATE_SERVER, server)
	}
	return 0
}

func (cache *RouteCache) RemoveServer(servicename string, key string) int {

	cache.Lock()
	defer cache.Unlock()
	if _, ret := cache.servers[servicename]; ret {
		for i, server := range cache.servers[servicename] {
			if server.Key == key {
				cache.servers[servicename] = append(cache.servers[servicename][:i], cache.servers[servicename][i+1:]...)
				if len(cache.servers[servicename]) == 0 {
					delete(cache.servers, servicename)
				}
				cache.circle[servicename].Remove(key)
				if len(cache.circle[servicename].Members()) == 0 {
					delete(cache.circle, servicename)
				}
				cache.handler.OnRouteCacheChangedHandleFunc(RELEASE_SERVER, server)
				return 0
			}
		}
	}
	return -1
}

func (cache *RouteCache) ClearServer() {

	cache.Lock()
	defer cache.Unlock()
	for servicename := range cache.servers {
		cache.servers[servicename] = cache.servers[servicename][0:0]
		delete(cache.servers, servicename)
		delete(cache.circle, servicename)
	}
}
