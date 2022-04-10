package loader

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ZentriaMC/apply-config/internal/core"
)

func JSONLoader(reader io.Reader) (data FileData, err error) {
	var d core.MapSection
	if err = json.NewDecoder(reader).Decode(&d); err != nil {
		return
	}

	data = &jsonFileData{
		data: d,
	}

	return
}

type jsonFileData struct {
	data core.MapSection
}

func (p *jsonFileData) Apply(values map[string]interface{}) (err error) {
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

func (p *jsonFileData) Save(wr io.Writer) (err error) {
	enc := json.NewEncoder(wr)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")

	err = enc.Encode(p.data)
	return
}
