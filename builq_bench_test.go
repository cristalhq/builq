package builq_test

import (
	"testing"

	"github.com/cristalhq/builq"
)

func BenchmarkBuildSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := builq.New()
		q("SELECT %s FROM %s", "foo,bar", "table")
		q("WHERE id = %$", 123)
		query, args, err := q.Build()

		switch {
		case err != nil:
			b.Fatal(err)
		case query == "":
			b.Fatal("empty query")
		case len(args) == 0:
			b.Fatal("empty args")
		}
	}
}

func BenchmarkBuildManyArgs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := builq.New()
		q("SELECT %s FROM %s", "foo,bar", "table")
		q("WHERE id = %$", 123)
		q("AND user_id = %$ OR referral_id = %$", 123, 456)
		q("OR status = %$ OR name = %$", "admin", "nice")
		q("AND password LIKE %$ OR password = %$", "str0ng", "qwerty")
		q("GROUP BY %$", "foo")
		q("ORDER BY %$", "bar")
		q("LIMIT %$;", 123)
		query, args, err := q.Build()

		switch {
		case err != nil:
			b.Fatal(err)
		case query == "":
			b.Fatal("empty query")
		case len(args) == 0:
			b.Fatal("empty args")
		}
	}
}

func BenchmarkBuildCached(b *testing.B) {
	q := builq.New()
	q("SELECT %s FROM %s", "foo,bar", "table")
	q("WHERE id = %$", 123)
	q("AND user_id = %$ OR referral_id = %$", 123, 456)
	q("OR status = %$ OR name = %$", "admin", "nice")
	q("GROUP BY %$", "foo")
	q("ORDER BY %$", "bar")
	q("LIMIT %$;", 123)
	query, args, err := q.Build()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		switch {
		case err != nil:
			b.Fatal(err)
		case query == "":
			b.Fatal("empty query")
		case len(args) == 0:
			b.Fatal("empty args")
		}
	}
}

func BenchmarkBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var bb builq.Builder
		bb.Addf("SELECT %s FROM %s", "foo,bar", "table")
		bb.Addf("WHERE id = %$", 123)
		query, args, err := bb.Build()

		switch {
		case err != nil:
			b.Fatal(err)
		case query == "":
			b.Fatal("empty query")
		case len(args) == 0:
			b.Fatal("empty args")
		}
	}
}
