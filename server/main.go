package main

import (
	"github.com/LgThinh/video-ranking-service/conf"
	"github.com/LgThinh/video-ranking-service/pkg/router"
)

// @securityDefinitions.apikey Authorization
// @in                         header
// @name                       Authorization

// @securityDefinitions.apikey EntityID
// @in                         header
// @name                       x-entity-id

// @securityDefinitions.apikey VideoID
// @in                         header
// @name                       x-video-id
func main() {
	conf.LoadConfig()
	router.NewRouter()
}
