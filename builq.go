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

func NewPostgreSQL() *Builder { return &Builder{pp: &postgres{}} }
func NewMySQL() *Builder      { return &Builder{pp: &mysql{}} }

func (b *Builder) Appendf(format string, args ...any) *Builder {
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
		fmt.Fprint(s, a.pp.NextPlaceholder())
	default:
		panic(fmt.Sprintf("unsupported verb %c", v))
	}
}

type placeholderProvider interface {
	NextPlaceholder() string
}

type postgres struct {
	counter int
}

func (p *postgres) NextPlaceholder() string {
	p.counter++
	return "$" + strconv.Itoa(p.counter)
}

type mysql struct{}

func (*mysql) NextPlaceholder() string { return "?" }
