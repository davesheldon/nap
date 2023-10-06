package naputil_test

import (
	"testing"

	"github.com/davesheldon/nap/naputil"
)

func TestCloneMap(t *testing.T) {
	t.Run("simple string map clone", func(t *testing.T) {
		key := "key1"
		val := "val1"
		newVal := "val2"

		original := map[string]string{key: val}
		cloned := naputil.CloneMap(original)

		if original[key] != cloned[key] {
			t.Errorf("Expected cloned[%s]=%s, got %s", key, val, cloned[key])
		}

		cloned[key] = newVal

		if original[key] == cloned[key] {
			t.Errorf("Expected original[%s]=%s, got %s", key, val, cloned[key])
		}
	})
}
