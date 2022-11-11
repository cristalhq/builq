package builq

import (
	"errors"
	"testing"
)

func TestBuilder(t *testing.T) {
	t.Run("unsupported verb", func(t *testing.T) {
		var b Builder
		b.Addf("SELECT * FROM %v", "users").
			Addf("LIMIT 100;")

		if _, _, err := b.Build(); err == nil {
			t.Errorf("must be error")
		} else if want := "unsupported verb v"; err.Error() != want {
			t.Errorf("\nhave: %v\nwant: %v", err, want)
		}
	})

	t.Run("different placeholders", func(t *testing.T) {
		var b Builder
		b.Addf("WHERE foo = %$ AND bar = %?", 1, 2)

		if _, _, err := b.Build(); !errors.Is(err, errMixedPlaceholders) {
			t.Errorf("\nhave: %v\nwant: %v", err, errMixedPlaceholders)
		}
	})

	t.Run("different placeholders in slices", func(t *testing.T) {
		var b Builder
		b.Addf("WHERE foo = %+$ AND bar = %+?", 1, 2)

		if _, _, err := b.Build(); !errors.Is(err, errNonSliceArgument) {
			t.Errorf("\nhave: %v\nwant: %v", err, errNonSliceArgument)
		}
	})

	t.Run("batch placeholder for non-slice args", func(t *testing.T) {
		var b Builder
		b.Addf("WHERE foo = %#$", 123)

		if _, _, err := b.Build(); !errors.Is(err, errNonSliceArgument) {
			t.Errorf("\nhave: %v\nwant: %v", err, errNonSliceArgument)
		}
	})
}
