package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"runtime"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/tidwall/jsonc"
	"github.com/urfave/cli/v2"
)

var currentCtx *ctx = nil

var templateFuncs = template.FuncMap{
	"hasvar": func(variable string) (ok bool) {
		_, ok = currentCtx.Data[variable]
		return
	},
	"getvar": func(variable string) (value string) {
		var ok bool
		value, ok = currentCtx.Data[variable]
		if !ok {
			panic(fmt.Errorf("undefined variable '%s'", variable))
		}
		return
	},
	"getvaror": func(variable string, def string) (value string) {
		var ok bool
		value, ok = currentCtx.Data[variable]
		if !ok {
			value = def
		}
		return
	},
	"semvercmp": func(variable, op, value string) bool {
		varValue, ok := currentCtx.Data[variable]
		if !ok {
			panic(fmt.Errorf("undefined variable '%s'", variable))
		}

		v1 := semver.MustParse(varValue)
		v2 := semver.MustParse(value)

		switch op {
		case "lt":
			return v1.LessThan(v2)
		case "le":
			return v1.LessThan(v2) || v1.Equal(v2)
		case "gt":
			return v1.GreaterThan(v2)
		case "ge":
			return v1.GreaterThan(v2) || v1.Equal(v2)
		case "eq":
			return v1.Equal(v2)
		default:
			panic(fmt.Errorf("unknown comparision op '%s'", op))
		}
	},
	"max": func(a, b int) int {
		return int(math.Max(float64(a), float64(b)))
	},
	"min": func(a, b int) int {
		return int(math.Min(float64(a), float64(b)))
	},
	"divide": func(a, b int) int {
		return a / b
	},
	"multiply": func(a, b int) int {
		return a * b
	},
}

func templateConfig(cctx *cli.Context, configBuf []byte) (buf []byte, err error) {
	dataRaw := cctx.StringSlice("data")
	dataFile := cctx.Path("vars-from")
	data := map[string]string{}

	if dataFile != "" {
		var varsBuf []byte
		if varsBuf, err = ioutil.ReadFile(dataFile); err != nil {
			err = fmt.Errorf("unable to read '%s': %w", dataFile, err)
			return
		}

		if err = json.Unmarshal(jsonc.ToJSON(varsBuf), &data); err != nil {
			err = fmt.Errorf("unable to parse '%s': %w", dataFile, err)
			return
		}
	}

	for _, pair := range dataRaw {
		split := strings.SplitN(pair, "=", 2)
		if len(split) != 2 {
			err = fmt.Errorf("invalid data flag '%s'\n", pair)
			return
		}
		data[split[0]] = split[1]
	}

	var tmpl *template.Template
	if tmpl, err = template.New("config").Funcs(templateFuncs).Parse(string(configBuf)); err != nil {
		err = fmt.Errorf("failed to template configuration: %w", err)
		return
	}

	currentCtx = &ctx{
		Data: data,
		Host: hostCtx{
			NumCPU: runtime.NumCPU(),
		},
	}

	var w bytes.Buffer
	if err = tmpl.Execute(&w, currentCtx); err != nil {
		err = fmt.Errorf("failed to template configuration: %w", err)
		return
	}

	buf = w.Bytes()
	return
}
