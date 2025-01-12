package main

import (
	"controller/api"
	"controller/utils/k8s"
	"flag"
	"log"
)

func main() {
	var dev bool

	flag.BoolVar(&dev, "dev", false, "Run the controller in development mode. This allows the controller to function OUTSIDE the K8s cluster.")
	flag.Parse()

	if err := k8s.Init(dev); err != nil {
		log.Panic(err)
	}
	if err := api.StartAPI(dev); err != nil {
		log.Panic(err)
	}
}
