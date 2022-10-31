package builq_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/cristalhq/builq"
)

func ExampleQuery1() {
	var b builq.Builder
	b.Appendf("SELECT %s FROM %s\n", "foo, bar", "users")
	b.Appendf("WHERE\n")
	b.Appendf("active IS TRUE\n")
	b.Appendf("AND user_id = %a OR invited_by = %a", 42, 42)
	query, _, _ := b.Build()

	fmt.Println(query)
	// Output:
	//
	// SELECT foo, bar FROM users
	// WHERE
	// active IS TRUE
	// AND user_id = $ OR invited_by = $
}

func ExampleQuery2() {
	var b builq.Builder
	b.Appendf("SELECT %s FROM %s\n", "foo, bar", "users")
	b.Appendf("WHERE\n")
	b.Appendf("active = %a\n", true)
	b.Appendf("AND user_id = %a\n", 42)
	b.Appendf("ORDER BY created_at\n")
	b.Appendf("LIMIT 100;")
	query, _, _ := b.Build()

	fmt.Println(query)
	// Output:
	//
	// SELECT foo, bar FROM users
	// WHERE
	// active = $
	// AND user_id = $
	// ORDER BY created_at
	// LIMIT 100;
}

func ExampleQuery3() {
	var b builq.Builder
	b.Appendf("SELECT * FROM foo\n")
	b.Appendf("WHERE active IS TRUE\n")
	b.Appendf("AND user_id = %a\n", 42)
	b.Appendf("LIMIT 100;")
	query, _, _ := b.Build()

	fmt.Println(query)
	// Output:
	//
	// SELECT * FROM foo
	// WHERE active IS TRUE
	// AND user_id = $
	// LIMIT 100;
}

func ExampleQuery4() {
	args := []interface{}{42, time.Now(), "just testing"}

	var b builq.Builder
	b.Appendf("INSERT (%s) INTO %s\n", getColumns(), "table")
	b.Appendf("VALUES (%a, %a, %a);", args...) // TODO(junk1tm): should %a support slices?
	query, _, _ := b.Build()

	fmt.Println(query)
	// Output:
	//
	// INSERT (id, created_at, value) INTO table
	// VALUES ($, $, $);
}

func getColumns() string {
	return strings.Join([]string{"id", "created_at", "value"}, ", ")
}
