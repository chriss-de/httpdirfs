package dirlist

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type HttpDirFsListingHTML struct{}

type HttpDirFsListingHTMLResult struct {
	dirContents []fs.DirEntry
	output      *bytes.Buffer
}

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)

func NewHtmlDirectoryListing() (hdfl *HttpDirFsListingHTML) {
	hdfl = &HttpDirFsListingHTML{}
	return hdfl
}

func (hdfl *HttpDirFsListingHTML) List(dirName string) (fd http.File, err error) {
	dirName = strings.TrimSuffix(dirName, "/index.html") + "/" // ???

	result := &HttpDirFsListingHTMLResult{}

	if result.dirContents, err = os.ReadDir(dirName); err != nil {
		return nil, err
	}

	result.output = new(bytes.Buffer)
	result.output.WriteString("<!doctype html>\n")
	result.output.WriteString("<meta name=\"viewport\" content=\"width=device-width\">\n")
	result.output.WriteString("<pre>\n")
	for _, ndc := range result.dirContents {
		name := ndc.Name()
		if ndc.IsDir() {
			name += "/"
		}
		u := url.URL{Path: name}
		result.output.WriteString(fmt.Sprintf("<a href=\"%s\" %s</a>\n", u.String(), htmlReplacer.Replace(name)))
	}
	result.output.WriteString("</pre>\n")

	return result, nil
}

func (r *HttpDirFsListingHTMLResult) Read(dst []byte) (int, error) {
	return r.output.Read(dst)
}
func (r *HttpDirFsListingHTMLResult) Seek(offset int64, whence int) (int64, error) {
	result := bytes.NewReader(r.output.Bytes())
	return result.Seek(offset, whence)
}
func (r *HttpDirFsListingHTMLResult) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, os.ErrInvalid
}
func (r *HttpDirFsListingHTMLResult) Stat() (os.FileInfo, error) {
	return r, nil
}
func (r *HttpDirFsListingHTMLResult) Close() error { return nil }
func (r *HttpDirFsListingHTMLResult) Name() string { return "listing.html" }
func (r *HttpDirFsListingHTMLResult) Size() int64  { return int64(r.output.Len()) }
func (r *HttpDirFsListingHTMLResult) Mode() fs.FileMode {
	return 0444
}
func (r *HttpDirFsListingHTMLResult) ModTime() time.Time { return time.Now() }
func (r *HttpDirFsListingHTMLResult) IsDir() bool        { return false }
func (r *HttpDirFsListingHTMLResult) Sys() any           { return nil }
