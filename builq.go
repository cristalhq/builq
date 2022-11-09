package builq

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

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

func (b *Builder) writeCountedArgs(s fmt.State, placeholder rune, args ...any) {
	for i, arg := range args {
		if i > 0 {
			fmt.Fprint(s, ", ")
		}
		b.appendArg(arg, placeholder)
		fmt.Fprintf(s, "$%d", b.counter)
	}
}

func (b *Builder) writeSimpleArgs(s fmt.State, placeholder rune, args ...any) {
	for i, arg := range args {
		if i > 0 {
			fmt.Fprint(s, ", ")
		}
		b.appendArg(arg, placeholder)
		fmt.Fprint(s, "?")
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

// argument is a wrapper for arguments passed to Builder.
type argument struct {
	value   any
	builder *Builder
}

// Format implements the [fmt.Formatter] interface.
func (a *argument) Format(s fmt.State, v rune) {
	var ok bool
	args := []any{a.value}

	switch v {
	case 's': // just a string
		fmt.Fprint(s, a.value)

	case '$': // PostgreSQL
		if s.Flag('+') {
			args, ok = a.asSlice(a.value)
			if !ok {
				return
			}
		}
		a.builder.writeCountedArgs(s, v, args...)

	case '?': // MySQL/SQLite
		if s.Flag('+') {
			args, ok = a.asSlice(a.value)
			if !ok {
				return
			}
		}
		a.builder.writeSimpleArgs(s, v, args...)

	default:
		a.builder.err = fmt.Errorf("unsupported verb %c", v)
	}
}

func (a *argument) asSlice(v any) ([]any, bool) {
	value := reflect.ValueOf(v)
	isSlice := value.Kind() == reflect.Slice

	if !isSlice {
		if a.builder.err == nil {
			a.builder.err = errNonSliceArgument
		}
		return nil, false
	}

	res := make([]any, value.Len())
	for i := 0; i < value.Len(); i++ {
		res[i] = value.Index(i).Interface()
	}
	return res, isSlice
}
