package main

import (
	"os"
	"strings"
	"sync"
)

type Configs map[string]string

var configs Configs
var once sync.Once

func init() {

}

func getConfigs() Configs {
	once.Do(loadConfigs)
	return configs
}

func (c Configs) get(key, defaultValue string) string {
	if value, found := c[key]; found {
		return value
	}
	return defaultValue

}

func loadConfigs() {
	configs = make(Configs)
	i := 1
	for ; i < len(os.Args); i++ {
		if key, found := strings.CutPrefix(os.Args[i], "--"); found {
			i++
			value := os.Args[i]
			configs[key] = value
		}
	}
}
