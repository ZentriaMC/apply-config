package loader

import "io"

type Loader func(io.Reader) (FileData, error)

type Configuration map[string](map[string]interface{})

type FileData interface {
	Apply(map[string]interface{}) error
	Save(io.Writer) error
}
