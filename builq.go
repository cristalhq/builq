package builq

import (
	"errors"
	"reflect"
	"strings"
)

// used to enforce const strings in API.
type constString string

// Columns is a convenience wrapper for table columns.
type Columns []string

// String implements the [fmt.Stringer] interface.
func (c Columns) String() string {
	return strings.Join(c, ", ")
}

// Prefixed acts the same as String but also prefixes each column with p.
func (c Columns) Prefixed(p string) string {
	return p + strings.Join(c, ", "+p)
}

// Builder for SQL queries.
type Builder struct {
	query       strings.Builder
	parts       []string
	rawArgs     [][]any
	args        []any
	err         error // the first error occurred while building the query.
	counter     int   // a counter for numbered placeholders ($1, $2, ...).
	placeholder rune  // a placeholder used to build the query.
	sep         byte
	debug       bool
}

// OnelineBuilder behaves like Builder but result is 1 line.
type OnelineBuilder struct {
	Builder
}

// Addf formats according to a format specifier, writes to query and appends args.
// Format param must be a constant string.
func (b *OnelineBuilder) Addf(format constString, args ...any) *Builder {
	if b.sep == 0 {
		b.sep = ' '
	}
	return b.addf(format, args...)
}

// Addf formats according to a format specifier, writes to query and appends args.
// Format param must be a constant string.
func (b *Builder) Addf(format constString, args ...any) *Builder {
	if b.sep == 0 {
		b.sep = '\n'
	}
	return b.addf(format, args...)
}

func (b *Builder) addf(format constString, args ...any) *Builder {
	if len(b.parts) == 0 {
		// TODO: better defaults
		b.parts = make([]string, 0, 10)
		b.rawArgs = make([][]any, 0, 10)
	}
	b.parts = append(b.parts, string(format))
	b.rawArgs = append(b.rawArgs, args)
	return b
}

func (b *Builder) Build() (query string, args []any, err error) {
	query = b.build()
	return query, b.args, b.err
}

func (b *Builder) DebugBuild() (query string) {
	b.debug = true
	q := b.build()
	b.debug = false
	return q
}

func (b *Builder) build() string {
	b.query.Grow(100)
	b.args = make([]any, 0, 10)

	for i := range b.parts {
		format := b.parts[i]
		args := b.rawArgs[i]

		err := b.write(format, args...)
		b.setErr(err)
	}
	return b.query.String()
}

func (b *Builder) asSlice(v any) []any {
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Slice {
		b.setErr(errNonSliceArgument)
		return nil
	}

	res := make([]any, value.Len())
	for i := 0; i < value.Len(); i++ {
		res[i] = value.Index(i).Interface()
	}
	return res
}

var (
	// errTooFewArguments is returned when too few arguments are provided to the
	// [Builder.Addf] method.
	errTooFewArguments = errors.New("too few arguments")

	// errTooFewArguments is returned when too many arguments are provided to
	// the [Builder.Addf] method.
	errTooManyArguments = errors.New("too many arguments")

	// errUnsupportedVerb is returned when an unsupported verb is found in a
	// query.
	errUnsupportedVerb = errors.New("unsupported verb")

	// errMixedPlaceholders is returned when different placeholders are used in
	// a single query (e.g. WHERE foo = %$ AND bar = %?).
	errMixedPlaceholders = errors.New("mixed placeholders must not be used in a single query")

	// errNonSliceArgument is returned when a non-slice argument is provided
	// with either `+` or `#` modifier.
	errNonSliceArgument = errors.New("non-slice arguments must not be used with slice modifiers")
)

func (b *Builder) setErr(err error) {
	if b.err == nil {
		b.err = err
	}
}
