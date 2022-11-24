//go:build prod
// +build prod

package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

var (

	//go:embed build/*
	build embed.FS
)

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

func EmbedHandler(prefix, root string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		defaultPath := fmt.Sprintf("%s/index.html", root)
		assetPath := path.Join(root, name)
		// If we can't find the asset, return the default index.html
		// build
		f, err := build.Open(assetPath)
		log.Info(assetPath, err)
		if os.IsNotExist(err) {
			return build.Open(defaultPath)
		}

		// Otherwise assume this is a legitimate request routed
		// correctly
		return f, err
	})

	return http.StripPrefix(prefix, http.FileServer(http.FS(handler)))
}

func init() {
	frontendHandler := EmbedHandler("/", "build")
	server.Handle("/", frontendHandler)
	server.Handle("/*filepath", frontendHandler)
}
