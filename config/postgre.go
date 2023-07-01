package config

import (
	"database/sql"
	"fmt"
	"github.com/aasumitro/tix/common"
	"log"
	"time"

	// postgresql
	_ "github.com/lib/pq"
)

func (cfg *Config) InitPostgresConn() {
	log.Println("Trying to open database connection pool . . . .")
	postgreSingleton.Do(func() {
		conn, err := sql.Open("postgres", cfg.PostgreDsnURL)
		if err != nil {
			panic(fmt.Sprintf("DATABASE_ERROR: %s", err.Error()))
		}
		conn.SetMaxOpenConns(common.DBMaxOpenConnection)
		conn.SetMaxIdleConns(common.DBMaxIdleConnection)
		conn.SetConnMaxLifetime(time.Duration(common.DBMaxLifetimeConnection))
		Postgre = conn
		if err := Postgre.Ping(); err != nil {
			panic(fmt.Sprintf("DATABASE_ERROR: %s", err.Error()))
		}
		log.Println("Database connection pool created with postgres driver . . . .")
	})
}
