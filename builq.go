package builq

import (
	"errors"
	"fmt"
	"strings"
)

// ErrMixedPlaceholders is returned by [Builder.Build] when different
// placeholders are used in a single query (e.g. WHERE foo = %$ AND bar = %?).
var ErrMixedPlaceholders = errors.New("mixed placeholders must not be used")

// Columns wrapper for your tables.
type Columns []string

func (c Columns) String() string {
	return strings.Join(c, ", ")
}

// Builder for SQL queries.
type Builder struct {
	query       strings.Builder
	args        []any
	err         error // the first error occurred while building the query.
	counter     int   // a counter for numbered placeholders ($1, $2, ...).
	placeholder rune  // a placeholder used to build the query.
}

// Addf formats according to a format specifier, writes to query and appends args.
func (b *Builder) Addf(format string, args ...any) *Builder {
	if b.err != nil {
		return b
	}

	wargs := make([]any, len(args))
	for i, arg := range args {
		wargs[i] = &argument{value: arg, builder: b}
	}

	fmt.Fprintf(&b.query, format+"\n", wargs...)

	return b
}

func (b *Builder) Build() (string, []any, error) {
	return b.query.String(), b.args, b.err
}

func (b *Builder) appendArg(arg any, placeholder rune) {
	if b.placeholder == 0 {
		b.placeholder = placeholder
	}
	if b.placeholder != placeholder {
		b.err = ErrMixedPlaceholders
	}
	b.args = append(b.args, arg)
}

// argument is a wrapper for arguments passed to Builder.
type argument struct {
	value   any
	builder *Builder
}

// Format implements the [fmt.Formatter] interface.
func (a *argument) Format(s fmt.State, v rune) {
	switch v {
	case 's': // just a string
		fmt.Fprint(s, a.value)

	case '$': // PostgreSQL
		a.builder.appendArg(a.value, v)
		a.builder.counter++
		fmt.Fprintf(s, "$%d", a.builder.counter)

	case '?': // MySQL/SQLite
		a.builder.appendArg(a.value, v)
		fmt.Fprint(s, "?")

	default:
		a.builder.err = fmt.Errorf("unsupported verb %c", v)
	}
}
