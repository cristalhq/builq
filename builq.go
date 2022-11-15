package builq

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

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

func (b *Builder) writeArgs(s fmt.State, verb rune, arg any, isMulti bool) {
	args := []any{arg}
	if isMulti {
		args = b.asSlice(arg)
	}

	for i, arg := range args {
		if i > 0 {
			fmt.Fprint(s, ", ")
		}

		switch verb {
		case '$': // PostgreSQL
			b.counter++
			fmt.Fprintf(s, "$%d", b.counter)
		case '?': // MySQL/SQLite
			fmt.Fprint(s, "?")
		default:
			panic("unreachable")
		}

		// store the first placeholder used in the query
		// to check for mixed placeholders later.
		if b.placeholder == 0 {
			b.placeholder = verb
		}
		if b.placeholder != verb {
			b.err = errMixedPlaceholders
			return
		}

		b.args = append(b.args, arg)
	}
}

func (b *Builder) writeBatchArgs(s fmt.State, verb rune, arg any) {
	args := b.asSlice(arg)
	for i, arg := range args {
		if i > 0 {
			fmt.Fprint(s, ", ")
		}
		fmt.Fprint(s, "(")
		b.writeArgs(s, verb, arg, true)
		fmt.Fprint(s, ")")
	}
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

var (
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

// argument is a wrapper for arguments passed to Builder.
type argument struct {
	value   any
	builder *Builder
}

// Format implements the [fmt.Formatter] interface.
func (a *argument) Format(s fmt.State, verb rune) {
	switch verb {
	case 's': // just a string
		fmt.Fprint(s, a.value)

	case '$', '?': // a query argument
		if s.Flag('#') {
			a.builder.writeBatchArgs(s, verb, a.value)
			return
		}
		isMulti := s.Flag('+')
		a.builder.writeArgs(s, verb, a.value, isMulti)

	default:
		a.builder.err = fmt.Errorf("%w %c", errUnsupportedVerb, verb)
	}
}
