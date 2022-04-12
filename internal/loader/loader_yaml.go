package loader

import (
	"fmt"
	"io"

	"github.com/ZentriaMC/apply-config/internal/core"
	"gopkg.in/yaml.v3"
)

func YAMLLoader(reader io.Reader) (data FileData, err error) {
	var d core.MapSection
	var buf []byte

	if buf, err = io.ReadAll(reader); err != nil {
		return
	}

	if err = yaml.Unmarshal(buf, &d); err != nil {
		return
	}

	if d == nil {
		d = core.MapSection{}
	}

	data = &yamlFileData{
		data: d,
	}

	return
}

type yamlFileData struct {
	data core.MapSection
}

func (p *yamlFileData) Apply(values map[string]interface{}) (err error) {
	var path []core.PathElement
	for k, v := range values {
		if path, err = core.ProcessPath(k); err != nil {
			err = fmt.Errorf("unable to parse key '%s': %w", k, err)
			return
		}

		_ = p.data.SetDeep(path, true, v)
	}

	return
}

func (p *yamlFileData) Save(wr io.Writer) (err error) {
	enc := yaml.NewEncoder(wr)
	enc.SetIndent(2)

	err = enc.Encode(p.data)
	return
}
