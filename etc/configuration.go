package etc

import "gopkg.in/yaml.v2"
import "github.com/gtlservice/gutils/logger"
import "github.com/gtlservice/gzkwrapper"

import (
	"io/ioutil"
	"os"
)

type Parameters map[string]string
type Cors Parameters

type Configuration struct {
	Version   string `yaml:"version,omitempty"`
	Fork      bool   `yaml:"fork,omitempty"`
	PidFile   string `yaml:"pidfile,omitempty"`
	ZkWrapper struct {
		Hosts     string `yaml:"hosts,omitempty"`
		Root      string `yaml:"root,omitempty"`
		Device    string `yaml:"device,omitempty"`
		Location  string `yaml:"location,omitempty"`
		OS        string `yaml:"os,omitempty"`
		Platform  string `yaml:"platform,omitempty"`
		Pulse     string `yaml:"pulse,omitempty"`
		Timeout   string `yaml:"timeout,omitempty"`
		Threshold int    `yaml:"threshold,omitempty"`
	} `yaml:"zkwrapper,omitempty"`
	HttpServer struct {
		Bind string `yaml:"bind,omitempty"`
		Cors Cors   `yaml:"cors,omitempty"`
	} `yaml:"http,omitempty"`
	Logger struct {
		LogFile  string `yaml:"logfile,omitempty"`
		LogLevel string `yaml:"loglevel,omitempty"`
		LogSize  int64  `yaml:"logsize,omitempty"`
	} `yaml:"logger,omitempty"`
}

var configuration *Configuration

func NewConfiguration(file string) (*Configuration, error) {

	fp, err := os.OpenFile(file, os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	c := makeDefault()
	if err := yaml.Unmarshal([]byte(data), c); err != nil {
		return nil, err
	}
	configuration = c
	return configuration, nil
}

func makeDefault() *Configuration {

	return &Configuration{
		Version: "default",
		Fork:    false,
		PidFile: "./gtlgateway.pid",
		ZkWrapper: struct {
			Hosts     string `yaml:"hosts,omitempty"`
			Root      string `yaml:"root,omitempty"`
			Device    string `yaml:"device,omitempty"`
			Location  string `yaml:"location,omitempty"`
			OS        string `yaml:"os,omitempty"`
			Platform  string `yaml:"platform,omitempty"`
			Pulse     string `yaml:"pulse,omitempty"`
			Timeout   string `yaml:"timeout,omitempty"`
			Threshold int    `yaml:"threshold,omitempty"`
		}{
			Hosts:     "127.0.0.1:2818",
			Root:      "/gtlservice",
			Device:    "",
			Location:  "center",
			OS:        "",
			Platform:  "",
			Pulse:     "10s",
			Timeout:   "30s",
			Threshold: 1,
		},
		HttpServer: struct {
			Bind string `yaml:"bind,omitempty"`
			Cors Cors   `yaml:"cors,omitempty"`
		}{
			Bind: ":8982",
			Cors: map[string]string{
				"origin":  "*",
				"methods": "GET",
			},
		},
		Logger: struct {
			LogFile  string `yaml:"logfile,omitempty"`
			LogLevel string `yaml:"loglevel,omitempty"`
			LogSize  int64  `yaml:"logsize,omitempty"`
		}{
			LogFile:  "logs/jobworker.log",
			LogLevel: "debug",
			LogSize:  2097152,
		},
	}
}

func GetConfiguration() *Configuration {

	return configuration
}

func (c *Configuration) GetVersion() string {

	if c != nil {
		return c.Version
	}
	return "default"
}

func (c *Configuration) GetFork() bool {

	if c != nil {
		return c.Fork
	}
	return false
}

func (c *Configuration) GetPidFile() string {

	if c != nil {
		return c.PidFile
	}
	return ""
}

func (c *Configuration) GetZkWrapper() *gzkwrapper.ServerArgs {

	if configuration != nil {
		return &gzkwrapper.ServerArgs{
			Hosts:     c.ZkWrapper.Hosts,
			Root:      c.ZkWrapper.Root,
			Device:    c.ZkWrapper.Device,
			Location:  c.ZkWrapper.Location,
			OS:        c.ZkWrapper.OS,
			Platform:  c.ZkWrapper.Platform,
			Pulse:     c.ZkWrapper.Pulse,
			Timeout:   c.ZkWrapper.Timeout,
			Threshold: c.ZkWrapper.Threshold,
		}
	}
	return nil
}

func (c *Configuration) GetHttpBind() string {

	if c != nil {
		return c.HttpServer.Bind
	}
	return ":30000"
}

func (c *Configuration) GetHttpCorsOrigin() string {

	if c != nil {
		return c.HttpServer.Cors["origin"]
	}
	return "*"
}

func (c *Configuration) GetHttpCorsMethod() string {

	if c != nil {
		return c.HttpServer.Cors["methods"]
	}
	return "GET"
}

func (c *Configuration) GetLogger() *logger.Args {

	if c != nil {
		return &logger.Args{
			FileName: c.Logger.LogFile,
			Level:    c.Logger.LogLevel,
			MaxSize:  c.Logger.LogSize,
		}
	}
	return nil
}
