package builq

import (
	"fmt"
	"strings"
)

// Builder for SQL queries.
type Builder struct {
	query   strings.Builder
	args    []any
	counter int   // a counter for numbered placeholders ($1, $2, ...).
	err     error // the first error occurred while building the query.
}

func (b *Builder) Appendf(format string, args ...any) *Builder {
	if b.err != nil {
		return b // return earlier if there is already an error.
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
		a.builder.counter++
		fmt.Fprintf(s, "$%d", a.builder.counter)

	case '?': // MySQL/SQLite
		a.builder.args = append(a.builder.args, a.value)
		fmt.Fprint(s, "?")

	default:
		a.builder.err = fmt.Errorf("builq: unsupported verb %c", v)
	}
}
