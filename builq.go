package builq

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Columns wrapper for your tables.
type Columns []string

func (c Columns) String() string {
	return strings.Join(c, ", ")
}

func (c Columns) Prefixed(p string) string {
	return p + strings.Join(c, ", "+p)
}

var (
	// errMixedPlaceholders is returned by [Builder.Build] when different
	// placeholders are used in a single query (e.g. WHERE foo = %$ AND bar = %?).
	errMixedPlaceholders = errors.New("mixed placeholders must not be used")

	errNonSliceArgument = errors.New("cannot expand non-slice argument")
)

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

func (b *Builder) writeArgs(s fmt.State, placeholder rune, arg any, isMulti bool) {
	args := []any{arg}
	if isMulti {
		args = b.asSlice(arg)
	}

	for i, arg := range args {
		if i > 0 {
			fmt.Fprint(s, ", ")
		}

		b.appendArg(arg, placeholder)

		switch placeholder {
		case '$': // PostgreSQL
			fmt.Fprintf(s, "$%d", b.counter)
		case '?': // MySQL/SQLite
			fmt.Fprint(s, "?")
		default:
			panic("unreachable")
		}
	}
}

func (b *Builder) appendArg(arg any, placeholder rune) {
	b.counter++
	if b.placeholder == 0 {
		b.placeholder = placeholder
	}
	if b.placeholder != placeholder {
		b.err = errMixedPlaceholders
	}
	b.args = append(b.args, arg)
}

func (b *Builder) asSlice(v any) []any {
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Slice {
		if b.err == nil {
			b.err = errNonSliceArgument
		}
		return nil
	}

	res := make([]any, value.Len())
	for i := 0; i < value.Len(); i++ {
		res[i] = value.Index(i).Interface()
	}
	return res
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

	case '$', '?': // PostgreSQL or MySQL/SQLite
		isMulti := s.Flag('+')
		a.builder.writeArgs(s, v, a.value, isMulti)

	default:
		a.builder.err = fmt.Errorf("unsupported verb %c", v)
	}
}
