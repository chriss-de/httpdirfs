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
		//fs := http.FileServer(httpdirfs.NewHttpDirFs("/", httpdirfs.WithDirectoryListing(dirlist.NewHtmlDirectoryListing())))
		//fs := http.FileServer(httpdirfs.NewHttpDirFs("/", httpdirfs.WithDirectoryListing(dirlist.NewJsonDirectoryListing())))
		//fs := http.FileServer(httpdirfs.NewHttpDirFs("/", httpdirfs.WithDirectoryListing(&httpdirfs.DefaultGolangListing{})))
		fs := http.FileServer(httpdirfs.NewHttpDirFs("/"))
		fs.ServeHTTP(w, r)
	})
}
