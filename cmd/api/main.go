package main

import (
	config "github.com/CLoouis/image-uploader"
	"github.com/CLoouis/image-uploader/pkg/api"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err.Error())
	}

	if err := api.Start(cfg); err != nil {
		panic(err.Error())
	}
}
