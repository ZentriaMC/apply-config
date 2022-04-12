package core_test

import (
	"testing"

	"github.com/ZentriaMC/apply-config/internal/core"
)

func TestNavigableMap(t *testing.T) {
	nm := core.MapSection{
		"one": map[string]interface{}{
			"two": map[string]interface{}{
				"three": 42,
			},
		},
	}

	path, err := core.ProcessPath(`one.two.three`)
	if err != nil {
		t.Fatal("unable to process path", err)
	}

	value, ok := nm.ValueDeep(path)
	if !ok {
		t.Errorf("value was not present")
		return
	}

	if value != 42 {
		t.Errorf("42 != %d", value)
		return
	}
}

func TestNavigableMapCreate(t *testing.T) {
	nm := &core.MapSection{}

	path, err := core.ProcessPath(`one.two.three`)
	if err != nil {
		t.Fatal("unable to process path", err)
	}

	var expected uint = 42
	_ = nm.SetDeep(path, true, expected)

	if v, ok := nm.ValueDeep(path); !ok || v != expected {
		t.Fatalf("not ok or %d != 42", v)
	}
}
