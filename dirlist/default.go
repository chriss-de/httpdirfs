package dirlist

import "net/http"

type DefaultGolangListing struct{}

func (d *DefaultGolangListing) List(string) (http.File, error) { return nil, nil }
