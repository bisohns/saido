/*
Copyright © 2021 Bisohns

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/bisohns/saido/client"
	"github.com/gorilla/handlers"
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	port   string
	server = http.NewServeMux()

	//go:embed build/*
	build embed.FS
)

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

func EmbedHandler(prefix, root string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(root, name)

		// If we can't find the asset, return the default index.html
		// build
		f, err := build.Open(assetPath)
		if os.IsNotExist(err) {
			return build.Open(fmt.Sprintf("%s/index.html", root))
		}

		// Otherwise assume this is a legitimate request routed
		// correctly
		return f, err
	})

	return http.StripPrefix(prefix, http.FileServer(http.FS(handler)))
}

var apiCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Run saido dashboard on a PORT env variable, fallback to set argument",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		hosts := client.NewHostsController(cfg)

		frontendHandler := EmbedHandler("/", "build")
		server.Handle("/", frontendHandler)
		server.Handle("/metrics", hosts)
		server.Handle("/*filepath", frontendHandler)
		log.Info("listening on :", port)
		_, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal(err)
		}
		go hosts.Run()
		loggedRouters := handlers.LoggingHandler(os.Stdout, server)
		// Trigger browser open
		url := fmt.Sprintf("http://localhost:%s", port)
		browser.OpenURL(url)
		if err := http.ListenAndServe(":"+port, loggedRouters); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	apiCmd.Flags().StringVarP(&port, "port", "p", "3000", "Port to run application server on")
	rootCmd.AddCommand(apiCmd)
}
