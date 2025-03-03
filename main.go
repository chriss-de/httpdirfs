package httpdirfs

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	dirlist "github.com/chriss-de/httpdirfs/dirlist"
)

type HttpDirFs struct {
	http.FileSystem
	rootPath         string
	tryFiles         []string
	directoryListing HttpDirFsListing
}

type HttpDirFsListing interface {
	List(string) (http.File, error)
}

func NewHttpDirFs(rootPath string, opts ...func(hdf *HttpDirFs)) (hdf *HttpDirFs, err error) {
	absRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	hdf = &HttpDirFs{rootPath: absRootPath}

	for _, opt := range opts {
		opt(hdf)
	}

	return hdf, nil
}

func WithDirectoryListing(dl HttpDirFsListing) func(lfs *HttpDirFs) {
	return func(hdf *HttpDirFs) {
		hdf.directoryListing = dl
	}
}

func WithTryFile(file string) func(hdf *HttpDirFs) {
	return func(hdf *HttpDirFs) {
		hdf.tryFiles = append(hdf.tryFiles, file)
	}
}

func WithTryFiles(files ...string) func(hdf *HttpDirFs) {
	return func(hdf *HttpDirFs) {
		hdf.tryFiles = append(hdf.tryFiles, files...)
	}
}

func (hdf *HttpDirFs) Open(name string) (fd http.File, err error) {
	var listOfFiles = []string{name}
	listOfFiles = append(listOfFiles, hdf.tryFiles...)

	return hdf.tryOpen(listOfFiles...)
}

func (hdf *HttpDirFs) tryOpen(fileNames ...string) (fd http.File, err error) {
	for tryOpenIdx, filename := range fileNames {
		filename = filepath.Join(hdf.rootPath, filename)
		if filename, err = filepath.Abs(filename); err != nil {
			return nil, err
		}
		if strings.HasPrefix(filename, hdf.rootPath) {
			fd, err = os.Open(filename)
			switch {
			// golang net/http tries for /index.html if `fd` is a directory
			case os.IsNotExist(err) && strings.HasSuffix(filename, "/index.html") && hdf.directoryListing != nil:
				switch hdf.directoryListing.(type) {
				case *dirlist.DefaultGolangListing:
					return nil, err
				default:
					return hdf.directoryListing.List(filename)
				}
			case os.IsNotExist(err) && tryOpenIdx < len(fileNames)-1:
				continue
			default:
				return fd, err
			}
		}
	}

	return nil, os.ErrNotExist
}
