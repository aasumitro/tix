package config

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
)

func (cfg *Config) InitSentryConn() {
	log.Println("Trying to initialize crash reporting handler . . . .")
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDsnURL,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		panic(fmt.Sprintf("Sentry initialization failed: %v\n", err))
	}
	log.Println("Crash reporting set to sentry . . . .")
}
