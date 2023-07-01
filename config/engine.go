package config

import (
	"github.com/aasumitro/tix/pkg/http/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

func (cfg *Config) InitGinEngine() {
	log.Println("Trying to init engine . . . .")
	engineSingleton.Do(func() {
		gin.SetMode(func() string {
			if cfg.AppDebug {
				return gin.DebugMode
			}
			return gin.ReleaseMode
		}())
		Engine = gin.Default()
		Engine.Use(middleware.CORS())
		log.Printf("Gin Engine (%s) created  . . . .", gin.Version)
	})
}
