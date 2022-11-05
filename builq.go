package builq

import (
	"fmt"
	"strconv"
	"strings"
)

// Builder for SQL queries.
type Builder struct {
	query strings.Builder
	args  []any
	pp    placeholderProvider
	err   error // the first error occurred while building the query.
}

func NewIterBuilder(placeholder string) Builder {
	return Builder{
		pp: &iterProvider{p: placeholder},
	}
}

func NewStaticBuilder(placeholder string) Builder {
	return Builder{
		pp: &staticProvider{p: placeholder},
	}
}

func (b *Builder) Appendf(format string, args ...any) *Builder {
	if b.err != nil {
		return b // return earlier if there is already an error.
	}

	if b.pp == nil {
		b.pp = &iterProvider{p: "$"}
	}

	wargs := make([]any, len(args))
	for i, arg := range args {
		wargs[i] = &argument{value: arg, pp: b.pp}
	}

	fmt.Fprintf(&b.query, format+"\n", wargs...)

	for _, warg := range wargs {
		arg := warg.(*argument)
		if err := arg.err; err != nil {
			b.err = err
			break
		}
		if arg.forQuery {
			b.args = append(b.args, arg.value)
		}
	}

	return b
}

func (b *Builder) Build() (string, []any, error) {
	if err := b.err; err != nil {
		return "", nil, err
	}
	return b.query.String(), b.args, nil
}

// argument is a wrapper for Printf-style arguments that implements fmt.Formatter.
type argument struct {
	value    any
	forQuery bool                // is it a query argument?
	pp       placeholderProvider // the source of the next placeholder.
	err      error               // an error occurred during the Format call.
}

// Format implements the fmt.Formatter interface.
func (a *argument) Format(s fmt.State, v rune) {
	switch v {
	case 's':
		// just a normal string (a table, a column, etc.), write it as is.
		fmt.Fprint(s, a.value)
	case 'a':
		// a query argument, mark it and write a placeholder.
		a.forQuery = true
		fmt.Fprint(s, a.pp.Next())
	default:
		a.err = fmt.Errorf("builq: unsupported verb %c", v)
		// panic(a.err) // this panic will be caught and written to s by the fmt code.
	}
}

type placeholderProvider interface {
	Next() string
}

type iterProvider struct {
	iter int
	p    string
}

func (ip *iterProvider) Next() string {
	ip.iter++
	return ip.p + strconv.Itoa(ip.iter)
}

type staticProvider struct {
	p string
}

func (sp *staticProvider) Next() string { return sp.p }
