# apply-config

Mass configuration editor

## Configuration format

Configuration is in [JSONC](https://github.com/tidwall/jsonc) format, you can use comments and trailing commas.

```jsonc
{
    "path/to/file.yaml": {
        "key": "value",
	"nested.key": "value",
	"nested": {}, // will replace previous `nested` declaration
	"nested.key1": 420
    }
}
```

(essentially `map[string]map[string]interface{}`)

Will yield:

```yaml
key: value
nested:
  key1: 420
```

### Semantics

Keys will be applied in order they are declared in the configuration - so you can overwrite previous declarations easily.

### Configuration templating

You can template configuration using [Golang native templating](https://pkg.go.dev/text/template).

Variables can be passed in via `--data=k=v` form or using `--vars-from path/to/vars.jsonc`.

## Example usage

```shell
$ ./apply-config --base-dir . --check --keep-going --config ./examples/mc.jsonc --vars-from ./examples/mc-vars.json
```
