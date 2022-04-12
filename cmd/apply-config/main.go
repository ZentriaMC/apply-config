package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZentriaMC/apply-config/internal/core"
	"github.com/ZentriaMC/apply-config/internal/loader"
	"github.com/tidwall/jsonc"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "apply-config",
		Description: "apply-config",
		Usage:       "Mass configuration editor",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "config",
				Usage:    "Configuration file to read declarations from",
				Required: true,
			},
			&cli.PathFlag{
				Name:     "base-dir",
				Usage:    "Base directory for relative file locations",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "keep-going",
				Usage:    "Keep going when declaration applying fails",
				Required: false,
				Value:    false,
			},
			&cli.BoolFlag{
				Name:     "check",
				Usage:    "Do not do anything, only check the declarations",
				Required: false,
				Value:    false,
			},
			&cli.StringSliceFlag{
				Name:     "data",
				Usage:    "Template data in form --data=k=v",
				Required: false,
			},
		},
		Action: entrypoint,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func entrypoint(cctx *cli.Context) (err error) {
	configFile := cctx.Path("config")
	keepGoing := cctx.Bool("keep-going")
	check := cctx.Bool("check")

	var configBuf []byte
	if configBuf, err = ioutil.ReadFile(configFile); err != nil {
		err = fmt.Errorf("unable to read config file '%s': %w", configFile, err)
		return
	}

	// Process configuration as a template
	var templatedConfig []byte
	if templatedConfig, err = templateConfig(cctx, configBuf); err != nil {
		return
	}

	var cfg loader.Configuration
	if err = json.Unmarshal(jsonc.ToJSON(templatedConfig), &cfg); err != nil {
		err = fmt.Errorf("unable to parse configuation: %w", err)
		return
	}

	baseDir := cctx.Path("base-dir")
	var baseDirAbs string
	if baseDirAbs, err = filepath.Abs(baseDir); err != nil {
		err = fmt.Errorf("unable to resolve absolute path for '%s': %w", baseDir, err)
		return
	}

	for filename, data := range cfg {
		// Filter out empty data
		if len(data) == 0 {
			continue
		}

		// Filter out bad paths
		// Note that we allow symlinks to go outside the base path
		fpath := filepath.Join(baseDirAbs, filename)
		if rel, rerr := filepath.Rel(baseDirAbs, fpath); rerr != nil {
			err = fmt.Errorf("file '%s' path is invalid: %w", filename, rerr)
			if keepGoing {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			return
		} else if strings.Contains(rel, "..") {
			err = fmt.Errorf("file '%s' path is invalid: goes outside from base directory", filename)
			if keepGoing {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			return
		}

		ext := filepath.Ext(filename)
		if len(ext) < 2 {
			continue
		}

		loader, ok := loaders[ext[1:]]
		if !ok || loader == nil {
			continue
		}

		if check {
			fmt.Printf("Checking file '%s' declarations\n", filename)
			hadError := false
			for k := range data {
				if _, err = core.ProcessPath(k); err != nil {
					hadError = true
					err = fmt.Errorf("unable to parse path '%s': %w", k, err)
					if keepGoing {
						fmt.Fprintln(os.Stderr, err)
						continue
					}
					return
				}
			}

			if !hadError {
				fmt.Printf("OK\n")
			}
		} else {
			if err = processFile(fpath, loader, data); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					//fmt.Printf("File '%s' does not exist\n", filename)
					continue
				}

				err = fmt.Errorf("unable to process file '%s': %w", fpath, err)
				if keepGoing {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				return
			}

			fmt.Printf("Processed file '%s'\n", filename)
		}
	}

	return
}

func processFile(path string, ldr loader.Loader, changes map[string]interface{}) (err error) {
	path = filepath.Clean(path)
	tmpPath := filepath.Join(filepath.Dir(path), fmt.Sprintf(".%s.tmp", filepath.Base(path)))

	var origFile *os.File
	var origStat os.FileInfo
	var targetFile *os.File

	// Get original file mode
	if origStat, err = os.Stat(path); err != nil {
		err = fmt.Errorf("unable to stat '%s': %w", path, err)
		return
	}

	if origFile, err = os.OpenFile(path, os.O_RDONLY, 0); err != nil {
		err = fmt.Errorf("unable to open '%s': %w", path, err)
		return
	}

	defer func() { _ = origFile.Close() }()

	if targetFile, err = os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE, origStat.Mode()); err != nil {
		err = fmt.Errorf("unable to open '%s': %w", tmpPath, err)
		return
	}

	defer func() { _ = targetFile.Close() }()

	// Process file
	var data loader.FileData
	if data, err = ldr(origFile); err != nil {
		err = fmt.Errorf("failed to load '%s': %w", path, err)
		return
	}

	if err = data.Apply(changes); err != nil {
		err = fmt.Errorf("failed to apply changes to '%s': %w", path, err)
		return
	}

	if err = data.Save(targetFile); err != nil {
		err = fmt.Errorf("failed to save changes to '%s': %w", tmpPath, err)
		return
	}

	_ = targetFile.Close()

	if err = os.Rename(tmpPath, path); err != nil {
		err = fmt.Errorf("failed to rename '%s' to '%s': %w", tmpPath, path, err)
		return
	}

	return
}
