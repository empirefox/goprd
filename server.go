package main

import (
	"github.com/empirefox/wsh2s"
	"github.com/empirefox/wsh2s/config"
)

func main() {
	config, err := config.LoadFromXpsWithEnv()
	if err != nil {
		panic(err)
	}

	s, err := wsh2s.NewServer(config)
	if err != nil {
		panic(err)
	}

	err = s.Serve()
	if err != nil {
		panic(err)
	}
}
