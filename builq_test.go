package builq

import (
	"errors"
	"strings"
	"testing"
)

func TestColumns(t *testing.T) {
	cols := Columns{"foo", "bar", "baz"}

	t.Run("string", func(t *testing.T) {
		const want = "foo, bar, baz"
		if got := cols.String(); got != want {
			t.Errorf("got %q; want %q", got, want)
		}
	})

	t.Run("prefixed", func(t *testing.T) {
		const want = "x.foo, x.bar, x.baz"
		if got := cols.Prefixed("x."); got != want {
			t.Errorf("got %q; want %q", got, want)
		}
	})
}

func TestBuilder(t *testing.T) {
	test := func(name string, wantErr error, format string, args ...any) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			var b Builder
			b.Addf(format, args...)
			_, _, err := b.Build()
			if !errors.Is(err, wantErr) {
				t.Errorf("\nhave: %v\nwant: %v", err, wantErr)
			}
		})
	}

	test("unsupported verb", errUnsupportedVerb, "SELECT * FROM %v", "users")
	test("mixed placeholders", errMixedPlaceholders, "WHERE foo = %$ AND bar = %?", 1, 2)
	test("non-slice argument", errNonSliceArgument, "WHERE foo = %+$", 1)
	test("non-slice argument (batch)", errNonSliceArgument, "WHERE foo = %#$", 1)
}

func FuzzBuilder(f *testing.F) {
	f.Add("SELECT %s FROM %s", "*", "users")
	f.Add("SELECT * FROM %s WHERE name = %$", "users", "john")
	f.Add("SELECT * FROM users WHERE name = %$ AND surname = %$", "john", "doe")

	// queries that won't return an error (probably should):
	// f.Add("%$%%$%*", "0", "0")

	f.Fuzz(func(t *testing.T, format, arg1, arg2 string) {
		var valid int
		for _, verb := range []string{"%s", "%$", "%?"} {
			valid += strings.Count(format, verb)
		}
		if valid != 2 {
			t.Skip("format must have 2 valid verbs")
		}

		var b Builder
		b.Addf(format, arg1, arg2)
		query, args, err := b.Build()

		if err != nil {
			// those errors are expected, we're looking for something new.
			if !errors.Is(err, errUnsupportedVerb) &&
				!errors.Is(err, errNonSliceArgument) &&
				!errors.Is(err, errMixedPlaceholders) {
				t.Errorf("unexpected error: %v", err)
			}
			return
		}

		// NOTE(junk1tm): fmt panics are written in query, should we parse it?
		_ = query

		checkArgs := func(strCnt, argsCnt int) {
			if strings.Count(format, "%s") == strCnt && len(args) != argsCnt {
				t.Errorf("format with %d string verbs must be bundled with %d arguments", strCnt, argsCnt)
			}
		}

		checkArgs(0, 2)
		checkArgs(1, 1)
		checkArgs(2, 0)
	})
}
