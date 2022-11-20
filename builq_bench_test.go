package builq_test

import (
	"testing"

	"github.com/cristalhq/builq"
)

func BenchmarkBuildNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var bb builq.Builder
		bb.Addf("SELECT %s FROM %s", "foo,bar", "table")
		bb.Addf("WHERE id = %$", 123)
		query, args, err := bb.Build()

		switch {
		case err != nil:
			b.Fatal(err)
		case len(query) == 0:
			b.Fatal("empty query")
		case len(args) == 0:
			b.Fatal("empty args")
		}
	}
}

func BenchmarkBuildManyArgs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var bb builq.Builder
		bb.Addf("SELECT %s FROM %s", "foo,bar", "table")
		bb.Addf("WHERE id = %$", 123)
		bb.Addf("AND user_id = %$ OR referral_id = %$", 123, 456)
		bb.Addf("OR status = %$ OR name = %$", "admin", "nice")
		bb.Addf("AND password LIKE %$ OR password = %$", "str0ng", "qwerty")
		bb.Addf("GROUP BY %$", "foo")
		bb.Addf("ORDER BY %$", "bar")
		bb.Addf("LIMIT %$;", 123)
		query, args, err := bb.Build()

		switch {
		case err != nil:
			b.Fatal(err)
		case len(query) == 0:
			b.Fatal("empty query")
		case len(args) == 0:
			b.Fatal("empty args")
		}
	}
}

func BenchmarkBuildCached(b *testing.B) {
	var bb builq.Builder
	bb.Addf("SELECT %s FROM %s", "foo,bar", "table")
	bb.Addf("WHERE id = %$", 123)
	bb.Addf("AND user_id = %$ OR referral_id = %$", 123, 456)
	bb.Addf("OR status = %$ OR name = %$", "admin", "nice")
	bb.Addf("GROUP BY %$", "foo")
	bb.Addf("ORDER BY %$", "bar")
	bb.Addf("LIMIT %$;", 123)
	query, args, err := bb.Build()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		switch {
		case err != nil:
			b.Fatal(err)
		case len(query) == 0:
			b.Fatal("empty query")
		case len(args) == 0:
			b.Fatal("empty args")
		}
	}
}
