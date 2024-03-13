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
	parts       []string
	args        [][]any
	err         error // the first error occurred while building the query.
	counter     int   // a counter for numbered placeholders ($1, $2, ...).
	placeholder byte  // a placeholder used to build the query.
	sep         byte  // a separator between Addf calls.
	debug       bool  // is it DebugBuild call to fill with args.
}

// OnelineBuilder behaves like Builder but result is 1 line.
type OnelineBuilder struct {
	Builder
}

// Q is a handy helper. Works as [NewOnline] and [Build] in one call.
func Q(format constString, args ...any) (query string, resArgs []any, err error) {
	return NewOneline()(format, args...).Build()
}

// BuildFn represents [Builder.Addf]. Just for the easier BuilderFunc declaration.
type BuildFn func(format constString, args ...any) *Builder

// Build the query and arguments.
func (q BuildFn) Build() (query string, args []any, err error) {
	return q("").Build()
}

// DebugBuild the query, good for debugging but not for REAL usage.
func (q BuildFn) DebugBuild() string {
	return q("").DebugBuild()
}

// New returns a new query builder, same as [Builder].
func New() BuildFn {
	var b Builder
	return b.Addf
}

// New returns a new query builder, same as [OnelineBuilder].
func NewOneline() BuildFn {
	var b OnelineBuilder
	return b.Addf
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

// Build the query and arguments.
func (b *Builder) Build() (query string, args []any, err error) {
	query, args = b.build()
	return query, args, b.err
}

// DebugBuild the query, good for debugging but not for REAL usage.
func (b *Builder) DebugBuild() (query string) {
	b.debug = true
	query, _ = b.build()
	b.debug = false
	return query
}

func (b *Builder) addf(format constString, args ...any) *Builder {
	if len(b.parts) == 0 {
		// TODO: better defaults
		b.parts = make([]string, 0, 10)
		b.args = make([][]any, 0, 10)
	}
	b.parts = append(b.parts, string(format))
	b.args = append(b.args, args)
	return b
}

func (b *Builder) build() (_ string, _ []any) {
	var query strings.Builder
	// TODO: better default (sum of parts + est len of indexes)
	query.Grow(100)
	// TODO: better default (count b.args in addf?)
	resArgs := make([]any, 0, 10)

	for i := range b.parts {
		format := b.parts[i]
		args := b.args[i]

		err := b.write(&query, &resArgs, format, args...)
		b.setErr(err)
	}

	// drop last separator for clarity.
	q := query.String()
	if q[len(q)-1] == b.sep {
		q = q[:len(q)-1]
	}
	return q, resArgs
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
	// errTooFewArguments passed to [Builder.Addf] method.
	errTooFewArguments = errors.New("too few arguments")

	// errTooFewArguments passed to [Builder.Addf] method.
	errTooManyArguments = errors.New("too many arguments")

	// errUnsupportedVerb when %X is found and X isn't supported.
	errUnsupportedVerb = errors.New("unsupported verb")

	// errIncorrectVerb is passed like `%+`.
	errIncorrectVerb = errors.New("incorrect verb")

	// errMixedPlaceholders when $ AND ? are mixed in 1 query.
	errMixedPlaceholders = errors.New("mixed placeholders must not be used in a single query")

	// errNonSliceArgument when a non-slice argument passed to placeholder with `+` or `#`.
	errNonSliceArgument = errors.New("non-slice arguments must not be used with slice modifiers")

	// errNonNumericArg expected number for %d but got something else.
	errNonNumericArg = errors.New("expected numeric argument")
)

func (b *Builder) setErr(err error) {
	if b.err == nil {
		b.err = err
	}
}
