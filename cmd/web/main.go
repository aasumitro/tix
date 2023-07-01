package main

import (
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/config"
	"github.com/aasumitro/tix/docs"
	"github.com/aasumitro/tix/internal"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/spf13/viper"
	"log"
)

//	@version                     1.0
//	@title                       TIX by BAKODE
//	@description                 Event Ticketing Management System
//
//	@contact.name                BAKODE SUPPORT
//	@contact.url                 https://bakode.xyz/
//	@contact.email               support@bakode.xyz
//
//	@securityDefinitions.jwt     ApiKeyAuth
//	@in                          header
//	@name                        Authorization
//
//	@externalDocs.description    OpenAPI
//	@externalDocs.url            https://swagger.io/resources/open-api/

func main() {
	// set config file
	viper.SetConfigFile(".env")
	// LOAD APP ENV
	config.LoadEnv()
	// INIT LANDLORD/MAIN DATABASE FOR DATA STORE
	config.Instance.InitPostgresConn()
	// INIT REDIS CONNECTION FOR DATA CACHE
	config.Instance.InitRedisConn()
	// INIT MAILER CONNECTION FOR EMAIL NOTIFICATION
	config.Instance.InitMailerConn()
	// INIT GIN ENGINE
	config.Instance.InitGinEngine()
	// INIT SENTRY
	if !config.Instance.AppDebug {
		config.Instance.InitSentryConn()
		config.Engine.Use(sentrygin.New(sentrygin.Options{}))
	}
	// INIT SWAGGER DOCS FEW CONFIG
	docs.SwaggerInfo.BasePath = config.Engine.BasePath()
	docs.SwaggerInfo.Host = config.Instance.AppURL
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	// INIT GOOGLE FORM SERVICE
	config.Instance.InitGoogleFormConn()

	// PRINT RUNNING LOG
	log.Printf("Run %s(%s)",
		config.Instance.AppName,
		common.Version)

	// BOOTSTRAP APP
	internal.RunApp(
		internal.WithEngine(config.Engine),
		internal.WithPostgreDatabase(config.Postgre),
		internal.WithRedisCache(config.Redis),
		internal.WithMailer(config.Mailer),
		internal.WithGoogleFormService(config.GoogleForm))

	// RUN SERVER
	log.Fatalln(config.Engine.Run(config.Instance.AppURL))
}
