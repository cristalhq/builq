package builq

import (
	"fmt"
	"strings"
)

// TODO(junk1tm): support PostgreSQL-style numbered placeholders.
var Placeholder = "$"

// Builder for SQL queries. The zero value is ready to use, just like strings.Builder.
type Builder struct {
	query strings.Builder
	args  []interface{}
}

func (b *Builder) Appendf(format string, args ...interface{}) *Builder {
	wargs := make([]interface{}, len(args))
	for i, arg := range args {
		wargs[i] = &argument{value: arg}
	}

	fmt.Fprintf(&b.query, format+"\n", wargs...)

	for _, warg := range wargs {
		if arg := warg.(*argument); arg.forQuery {
			b.args = append(b.args, arg.value)
		}
	}

	return b
}

func (b *Builder) Build() (string, []interface{}, error) {
	return b.query.String(), b.args, nil
}

// argument is a wrapper for Printf-style arguments that implements fmt.Formatter.
type argument struct {
	value    interface{}
	forQuery bool // is it a query argument?
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
		fmt.Fprint(s, Placeholder)
	default:
		panic(fmt.Sprintf("unsupported verb %c", v))
	}
}
