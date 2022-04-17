package main

import (
	"bytes"
	"io"

	"github.com/ZentriaMC/apply-config/internal/loader"
	"github.com/tidwall/jsonc"
)

var loaders = map[string]loader.Loader{
	"yaml":       loader.YAMLLoader,
	"yml":        loader.YAMLLoader,
	"properties": loader.JavaPropertiesLoader,
	"json":       loader.JSONLoader,
	"jsonc": func(r io.Reader) (fd loader.FileData, err error) {
		var buf bytes.Buffer
		if _, err = io.Copy(&buf, r); err != nil {
			return
		}

		jsoncBuf := jsonc.ToJSONInPlace(buf.Bytes())
		return loader.JSONLoader(bytes.NewReader(jsoncBuf))
	},
	//"xml": nil
}
