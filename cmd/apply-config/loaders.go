package main

import "github.com/ZentriaMC/apply-config/internal/loader"

var loaders = map[string]loader.Loader{
	"yaml":       loader.YAMLLoader,
	"yml":        loader.YAMLLoader,
	"properties": loader.JavaPropertiesLoader,
	"json":       loader.JSONLoader,
	//"xml": nil
}
