package builq

import (
	"fmt"
	"strings"
)

// Builder for SQL queries.
type Builder struct {
	query strings.Builder
	args  []interface{}
}

// Newf creates a new query
func Newf(format string, args ...interface{}) *Builder {
	b := &Builder{}
	q := fmt.Sprintf(format, args...)
	b.query.WriteString(q)
	return b
}

// Query that was build.
func (b *Builder) Query() string { return b.query.String() }

// Args passed to the query.
func (b *Builder) Args() []interface{} { return b.args }

// Append the expression with prepended "\n"
// Optional expressions are prepended with " ".
func (b *Builder) Append(expr string, exprs ...string) {
	b.query.WriteString("\n" + expr)
	for _, expr := range exprs {
		b.query.WriteByte(' ')
		b.query.WriteString(expr)
	}
}

// Add the expression and the parameter and return its index in query.
func (b *Builder) Add(expr string, v interface{}) string {
	param := b.AddParam(v)
	b.query.WriteString("\n" + expr + param)
	return param
}

// AddParam add parameter and return its index in query.
func (b *Builder) AddParam(v interface{}) string {
	b.args = append(b.args, v)
	return fmt.Sprintf("$%d", len(b.args))
}

// AddParams add parameters and return its index in query.
func (b *Builder) AddParams(v ...interface{}) string {
	res := ""
	for i, v := range v {
		if i > 0 {
			res += ", "
		}
		res += b.AddParam(v)
	}
	return res
}
