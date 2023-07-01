package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed all:build
var resource embed.FS

func SPAAssets() (spa fs.FS) {
	var err error
	if spa, err = fs.Sub(resource, "build"); err != nil {
		log.Fatalln("SPA_FS_ERR:", err.Error())
	}
	return spa
}
