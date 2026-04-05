package main

import (
	"os"
	"strings"
)

var configs map[string]string

func init() {
	loadConfigs()
}

func loadConfigs() {
	configs = make(map[string]string)
	i := 1
	for ; i < len(os.Args); i++ {
		if key, found := strings.CutPrefix(os.Args[i], "--"); found {
			i++
			value := os.Args[i]
			configs[key] = value
		}
	}
}
