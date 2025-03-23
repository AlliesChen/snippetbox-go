package ui

import (
	"embed"
	"io/fs"
)

//go:embed static/css static/js static/img html
var Files embed.FS
var StaticFS fs.FS

func init() {
	var err error
	StaticFS, err = fs.Sub(Files, "static")
	if err != nil {
		panic(err)
	}
}
