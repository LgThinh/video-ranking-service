package main

import (
	"github.com/LgThinh/video-ranking-service/conf"
	"github.com/LgThinh/video-ranking-service/pkg/router"
)

// @securityDefinitions.apikey Authorization
// @in                         header
// @name                       Authorization

// @securityDefinitions.apikey User ID
// @in                         header
// @name                       x-user-id
func main() {
	conf.LoadConfig()
	router.NewRouter()
}
