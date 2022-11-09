package builq_test

import (
	"errors"
	"testing"

	"github.com/cristalhq/builq"
)

func TestBuilder(t *testing.T) {
	t.Run("unsupported verb", func(t *testing.T) {
		var b builq.Builder
		b.Addf("SELECT * FROM %v", "users")
		if _, _, err := b.Build(); err == nil {
			t.Errorf("want an error")
		}
	})

	t.Run("different placeholders", func(t *testing.T) {
		var b builq.Builder
		b.Addf("WHERE foo = %$ AND bar = %?", 1, 2)
		if _, _, err := b.Build(); !errors.Is(err, builq.ErrMixedPlaceholders) {
			t.Errorf("got %v; want %v", err, builq.ErrMixedPlaceholders)
		}
	})
}
