package httpdirfs

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type HttpDirFsListingJSON struct{}
type HttpDirFsListingJSONResult struct {
	dirContents []fs.DirEntry
	output      *bytes.Buffer
}

func NewJsonDirectoryListing() (hdfl *HttpDirFsListingJSON) {
	hdfl = &HttpDirFsListingJSON{}
	return hdfl
}

func (list *HttpDirFsListingJSON) List(dirName string) (fd http.File, err error) {
	dirName = strings.TrimSuffix(dirName, "/index.html") + "/" // ???

	result := &HttpDirFsListingJSONResult{}

	if result.dirContents, err = os.ReadDir(dirName); err != nil {
		return nil, err
	}

	result.output = new(bytes.Buffer)
	result.output.WriteString(`{`)
	result.output.WriteString(`"fileCount":` + strconv.Itoa(len(result.dirContents)) + `,`)
	//result.output.WriteString(`"directory":"` + dirName + `",`)
	result.output.WriteString(`"files":[`)
	for idx, dc := range result.dirContents {
		name := dc.Name()
		fileType := "file"
		if dc.IsDir() {
			fileType = "directory"
		}
		dcInfo, err := dc.Info()
		if err != nil {
			return nil, err
		}
		result.output.WriteString(
			fmt.Sprintf(`{"name":"%s","type":"%s","mode":"%s","size":%d,"mtime":"%s"}`,
				name,
				fileType,
				dc.Type().String(),
				dcInfo.Size(),
				dcInfo.ModTime().Format(time.RFC3339),
			),
		)
		if idx < len(result.dirContents)-1 {
			result.output.WriteString(",")
		}
	}
	result.output.WriteString("]}")

	return result, nil
}

func (r *HttpDirFsListingJSONResult) Read(dst []byte) (int, error) {
	return r.output.Read(dst)
}
func (r *HttpDirFsListingJSONResult) Seek(offset int64, whence int) (int64, error) {
	result := bytes.NewReader(r.output.Bytes())
	return result.Seek(offset, whence)
}
func (r *HttpDirFsListingJSONResult) Readdir(int) ([]fs.FileInfo, error) {
	return nil, os.ErrInvalid
}
func (r *HttpDirFsListingJSONResult) Stat() (os.FileInfo, error) {
	return r, nil
}
func (r *HttpDirFsListingJSONResult) Close() error { return nil }
func (r *HttpDirFsListingJSONResult) Name() string { return "listing.json" }
func (r *HttpDirFsListingJSONResult) Size() int64  { return int64(r.output.Len()) }
func (r *HttpDirFsListingJSONResult) Mode() fs.FileMode {
	return 0444
}
func (r *HttpDirFsListingJSONResult) ModTime() time.Time { return time.Now() }
func (r *HttpDirFsListingJSONResult) IsDir() bool        { return false }
func (r *HttpDirFsListingJSONResult) Sys() any           { return nil }
