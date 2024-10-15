package main

import (
	"assigement_wallet/config"
	"assigement_wallet/http_core"
	"assigement_wallet/pkg/db_util"
	"assigement_wallet/pkg/redis_util"
	"assigement_wallet/router"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	if os.Getenv("ReleaseMode") == "true" {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(http_core.SessionControl)
	router.Register(g)
	if os.Getenv("PPROF") == "true" {
		pprof.Register(g, "/wallet/_pprof_/")
	}
	//init system
	err := config.InitConfigs()
	if err != nil {
		log.Fatal("config.InitConfigs:", err)
	}

	err = db_util.InitDB()
	if err != nil {
		log.Println("db_util.InitDB:", err)
	}
	err = redis_util.InitRedis()
	if err != nil {
		log.Fatal("redis_util.InitRedis:", err)
	}
	db_util.InitDemoData(db_util.GetDB())
	fmt.Println("server started...")
	g.Run(":80")

}
