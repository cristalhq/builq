package builq

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
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
	args        []any
	formats     []string
	rawArgs     [][]any
	debug       bool
	err         error // the first error occurred while building the query.
	counter     int   // a counter for numbered placeholders ($1, $2, ...).
	placeholder rune  // a placeholder used to build the query.
}

type constString string

// Addf formats according to a format specifier, writes to query and appends args.
// Format param must be a constant string.
func (b *Builder) Addf(format constString, args ...any) *Builder {
	if b.formats == nil {
		b.formats = make([]string, 0, 10)
		b.rawArgs = make([][]any, 0, 10)
	}
	b.formats = append(b.formats, string(format))
	b.rawArgs = append(b.rawArgs, args)
	return b
}

func (b *Builder) Build() (string, []any, error) {
	query := b.build()
	return query, b.args, b.err
}

func (b *Builder) DebugBuild() string {
	b.debug = true
	query := b.build()
	b.debug = false
	return query
}

func (b *Builder) build() string {
	sb := make([]string, 0, len(b.formats))

	for i := range b.formats {
		format := b.formats[i]
		args := b.rawArgs[i]

		argID := -1
		res := re.ReplaceAllStringFunc(format, func(s string) string {
			argID++

			if len(args) == 0 || len(args) == argID {
				b.setErr(errTooFewArguments)
				return ""
			}
			return b.write(s, args[argID])
		})

		if len(args) > 0 && argID != len(args)-1 {
			b.setErr(errTooManyArguments)
		}
		sb = append(sb, res, "\n")
	}
	return strings.Join(sb, "")
}

func (b *Builder) write(s string, arg any) string {
	var w strings.Builder

	if len(s) == 2 { // %s %$ %?
		b.writeArgs(&w, rune(s[1]), arg, false)
	} else { // %+$ %#$ %+? %#?
		verb := rune(s[2])
		isMulti := rune(s[1]) == '+'
		isBatch := rune(s[1]) == '#'
		if isBatch {
			b.writeBatchArgs(&w, verb, arg)
		} else {
			b.writeArgs(&w, verb, arg, isMulti)
		}
	}

	return w.String()
}

func (b *Builder) writeArgs(w io.Writer, verb rune, arg any, isMulti bool) {
	args := []any{arg}
	if isMulti {
		args = b.asSlice(arg)
	}

	for i, arg := range args {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		if b.debug {
			switch arg := arg.(type) {
			case string:
				fmt.Fprint(w, "'"+arg+"'")
			default:
				fmt.Fprint(w, arg)
			}
			continue
		}

		switch verb {
		case 's': // just a string
			fmt.Fprint(w, arg)
		case '$': // PostgreSQL
			b.counter++
			fmt.Fprintf(w, "$%d", b.counter)
			b.args = append(b.args, arg)
		case '?': // MySQL/SQLite
			fmt.Fprint(w, "?")
			b.args = append(b.args, arg)
		default:
			panic("unreachable")
		}

		// store the first placeholder used in the query
		// to check for mixed placeholders later
		// ignore 's' because it's ok to have in any case
		if b.placeholder == 0 && verb != 's' {
			b.placeholder = verb
		}
		if verb != 's' && b.placeholder != verb {
			b.setErr(errMixedPlaceholders)
			return
		}
	}
}

func (b *Builder) writeBatchArgs(s io.Writer, verb rune, arg any) {
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
		b.setErr(errNonSliceArgument)
		return nil
	}

	res := make([]any, value.Len())
	for i := 0; i < value.Len(); i++ {
		res[i] = value.Index(i).Interface()
	}
	return res
}

func (b *Builder) setErr(err error) {
	if b.err == nil {
		b.err = err
	}
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

// match any of: %s %$ %+$ %#$ %? %+? %#?
var re = regexp.MustCompile("%(s|(\\+|#|)\\$|(\\+|#|)\\?)")
