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
	if b.pp == nil {
		b.pp = &iterProvider{p: "$"}
	}

	wargs := make([]any, len(args))
	for i, arg := range args {
		wargs[i] = &argument{value: arg, pp: b.pp}
	}

	fmt.Fprintf(&b.query, format+"\n", wargs...)

	for _, warg := range wargs {
		if arg := warg.(*argument); arg.forQuery {
			b.args = append(b.args, arg.value)
		}
	}

	return b
}

func (b *Builder) Build() (string, []any, error) {
	return b.query.String(), b.args, nil
}

// argument is a wrapper for Printf-style arguments that implements fmt.Formatter.
type argument struct {
	value    any
	forQuery bool // is it a query argument?
	pp       placeholderProvider
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
		panic(fmt.Sprintf("unsupported verb %c", v))
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
