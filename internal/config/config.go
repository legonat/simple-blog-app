package config

import (
	"awesomeProjectRentaTeam/pkg/erx"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"sync"
)

type BlogDb struct {
	Path string
}

type DummyText struct {
	Path string
}

type Logger struct {
	Path string
}

type Server struct {
	Domain     string
	Port       int
	TlsEnabled bool
}

type Config struct {
	BlogDb BlogDb
	DummyText
	Logger Logger
	Server Server
}

var instance *Config
var once sync.Once

func GetConfigInstance() *Config {
	once.Do(func() {
		instance = &Config{}
		err := GetConfig(instance)
		if err != nil {
			panic(erx.New(err))
		}
	})
	return instance
}

func GetConfig(config *Config) error {
	file, err := ioutil.ReadFile("./app.conf")
	if err != nil {
		return erx.New(err)
	}
	err = toml.Unmarshal(file, config)
	if err != nil {
		return erx.New(err)
	}
	return nil
}
