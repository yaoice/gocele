package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yaoice/gocele/pkg/config"
	"github.com/yaoice/gocele/pkg/log"
	"github.com/yaoice/gocele/pkg/route"
	"github.com/yaoice/gocele/pkg/util"
)

func main() {
	utils.SetupDumpStacksTrap()
	r := gin.Default()
	m := config.GetString(config.FLAG_KEY_GIN_MODE)
	gin.SetMode(m)

	route.InstallRoutes(r)
	serverBindAddr := fmt.Sprintf("%s:%d", config.GetString(config.FLAG_KEY_SERVER_HOST), config.GetInt(config.FLAG_KEY_SERVER_PORT))
	log.Infof("Run server at %s", serverBindAddr)
	r.Run(serverBindAddr) // listen and serve
}
