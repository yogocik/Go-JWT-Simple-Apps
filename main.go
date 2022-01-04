package main

import (
	"integration/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	appConfig := config.NewConfig()
	defer func() {
		if err := appConfig.InfraManager.SqlDb(); err != nil {
			panic(err)
		}
	}()
	defer func() {
		if err := appConfig.InfraManager.KVStorage(); err != nil {
			panic(err)
		}
	}()
	routeEngine := appConfig.Routes
	err := routeEngine.RouterEngine.Run(appConfig.ApiBaseUrl)
	if err != nil {
		log.Fatal(err)
	}
}
