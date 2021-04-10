package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed webapp/build/*
var content embed.FS

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

func main() {
	router := gin.Default()

	router.Use(static.Serve("/", EmbedFolder(content, "webapp/build")))

	router.NoRoute(func(c *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	router.Run(":80")
}
