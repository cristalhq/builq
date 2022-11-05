package builq_test

import (
	"testing"

	"github.com/cristalhq/builq"
)

func TestBuilder(t *testing.T) {
	t.Run("unsupported verb", func(t *testing.T) {
		var b builq.Builder
		b.Appendf("SELECT * FROM %v", "users")
		if _, _, err := b.Build(); err == nil {
			t.Errorf("want an error")
		}
	})
}
