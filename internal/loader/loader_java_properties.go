package loader

import (
	"fmt"
	"io"

	"github.com/ZentriaMC/apply-config/internal/core"
	jproperties "github.com/magiconair/properties"
)

func JavaPropertiesLoader(reader io.Reader) (data FileData, err error) {
	var buf []byte
	var encoding jproperties.Encoding
	var properties *jproperties.Properties

	// TODO
	//nolint
	encoding = jproperties.UTF8

	if buf, err = io.ReadAll(reader); err != nil {
		return
	}

	if properties, err = jproperties.Load(buf, encoding); err != nil {
		return
	}

	data = &javaPropertiesFileData{
		properties: properties,
		encoding:   encoding,
	}

	return
}

type javaPropertiesFileData struct {
	properties *jproperties.Properties
	encoding   jproperties.Encoding
}

func (p *javaPropertiesFileData) Apply(values map[string]interface{}) (err error) {
	var path []core.PathElement
	for k, v := range values {
		if path, err = core.ProcessPath(k); err != nil {
			err = fmt.Errorf("unable to parse key '%s': %w", k, err)
			return
		}

		if len(path) != 1 {
			err = fmt.Errorf("key '%s' is invalid for Java Properties", k)
			return
		}

		first := path[0]
		if _, ok := first.(*core.ObjectPathElement); !ok {
			err = fmt.Errorf("key '%s' is invalid for Java Properties", k)
			return
		}

		if v != nil {
			if err = p.properties.SetValue(first.String(), v); err != nil {
				return
			}
		} else {
			p.properties.Delete(first.String())
		}
	}
	return
}

func (p *javaPropertiesFileData) Save(wr io.Writer) (err error) {
	_, err = p.properties.Write(wr, p.encoding)
	return
}
