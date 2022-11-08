package builq

import (
	"errors"
	"fmt"
	"strings"
)

// ErrDifferentPlaceholders is returned by [Builder.Build] when different
// placeholders are used in a single query (e.g. WHERE foo = %$ AND bar = %?).
var ErrDifferentPlaceholders = errors.New("builq: different placeholders must not be used")

// Builder for SQL queries.
type Builder struct {
	query        strings.Builder
	args         []any
	counter      int               // a counter for numbered placeholders ($1, $2, ...).
	err          error             // the first error occurred while building the query.
	placeholders map[rune]struct{} // a set of placeholders used to build the query.
}

func (b *Builder) Appendf(format string, args ...any) *Builder {
	if b.placeholders == nil {
		b.placeholders = make(map[rune]struct{})
	}

	// return earlier if there is already an error.
	if b.err != nil {
		return b
	}

	wargs := make([]any, len(args))
	for i, arg := range args {
		wargs[i] = &argument{value: arg, builder: b}
	}

	// writing to strings.Builder always returns no error.
	_, _ = fmt.Fprintf(&b.query, format+"\n", wargs...)

	return b
}

func (b *Builder) Build() (string, []any, error) {
	if len(b.placeholders) > 1 {
		return "", nil, ErrDifferentPlaceholders
	}
	return b.query.String(), b.args, b.err
}

// argument is a wrapper for arguments passed to [Builder.Appendf].
type argument struct {
	value   any
	builder *Builder
}

// Format implements the [fmt.Formatter] interface.
func (a *argument) Format(s fmt.State, v rune) {
	switch v {
	case 's': // table/column/etc.
		fmt.Fprint(s, a.value)

	case '$': // PostgreSQL
		a.builder.args = append(a.builder.args, a.value)
		a.builder.placeholders[v] = struct{}{}
		a.builder.counter++
		fmt.Fprintf(s, "$%d", a.builder.counter)

	case '?': // MySQL/SQLite
		a.builder.args = append(a.builder.args, a.value)
		a.builder.placeholders[v] = struct{}{}
		fmt.Fprint(s, "?")

	default:
		a.builder.err = fmt.Errorf("builq: unsupported verb %c", v)
	}
}
