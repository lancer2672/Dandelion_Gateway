package main

import (
	"github.com/lancer2672/Dandelion_Gateway/internal/helper"
	"github.com/lancer2672/Dandelion_Gateway/internal/utils"
	"github.com/lancer2672/Dandelion_Gateway/server"
	"github.com/lancer2672/Dandelion_Gateway/services"
)

func main() {
	utils.LoadConfig(".")
	helper.ConfigHttpClient()
	services.ConfigServices()
	server.RunServer()
}
