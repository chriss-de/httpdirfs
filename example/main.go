package main

import (
	"net/http"

	"github.com/chriss-de/httpdirfs"
)

func main() {
	httpServer := &http.Server{Handler: serve(), Addr: ":9080"}
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func serve() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//httpDir, err := httpdirfs.NewHttpDirFs("/", httpdirfs.WithDirectoryListing(dirlist.NewHtmlDirectoryListing()))
		//httpDir, err := httpdirfs.NewHttpDirFs("/", httpdirfs.WithDirectoryListing(dirlist.NewJsonDirectoryListing()))
		//httpDir, err := httpdirfs.NewHttpDirFs("/", httpdirfs.WithDirectoryListing(&dirlist.DefaultGolangListing{}))
		httpDir, err := httpdirfs.NewHttpDirFs("/")
		if err != nil {
			panic(err)
		}

		fs := http.FileServer(httpDir)
		fs.ServeHTTP(w, r)
	})
}
