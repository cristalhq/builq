package builq_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/cristalhq/builq"
)

func ExampleQuery1() {
	b := builq.NewPostgreSQL()
	b.Appendf("SELECT %s FROM %s", "foo, bar", "users")
	b.Appendf("WHERE")
	b.Appendf("active IS TRUE")
	b.Appendf("AND user_id = %a OR invited_by = %a", 42, 42)
	query, _, _ := b.Build()

	fmt.Println(query)
	// Output:
	//
	// SELECT foo, bar FROM users
	// WHERE
	// active IS TRUE
	// AND user_id = $1 OR invited_by = $2
}

func ExampleQuery2() {
	b := builq.NewPostgreSQL()
	b.Appendf("SELECT %s FROM %s", "foo, bar", "users")
	b.Appendf("WHERE")
	b.Appendf("active = %a", true)
	b.Appendf("AND user_id = %a", 42)
	b.Appendf("ORDER BY created_at")
	b.Appendf("LIMIT 100;")
	query, _, _ := b.Build()

	fmt.Println(query)
	// Output:
	//
	// SELECT foo, bar FROM users
	// WHERE
	// active = $1
	// AND user_id = $2
	// ORDER BY created_at
	// LIMIT 100;
}

func ExampleQuery3() {
	b := builq.NewPostgreSQL()
	b.Appendf("SELECT * FROM foo")
	b.Appendf("WHERE active IS TRUE")
	b.Appendf("AND user_id = %a", 42)
	b.Appendf("LIMIT 100;")
	query, _, _ := b.Build()

	fmt.Println(query)
	// Output:
	//
	// SELECT * FROM foo
	// WHERE active IS TRUE
	// AND user_id = $1
	// LIMIT 100;
}

func ExampleQuery4() {
	args := []any{42, time.Now(), "just testing"}

	b := builq.NewMySQL()
	b.Appendf("INSERT (%s) INTO %s", getColumns(), "table")
	b.Appendf("VALUES (%a, %a, %a);", args...) // TODO(junk1tm): should %a support slices?
	query, _, _ := b.Build()

	fmt.Println(query)
	// Output:
	//
	// INSERT (id, created_at, value) INTO table
	// VALUES (?, ?, ?);
}

func getColumns() string {
	return strings.Join([]string{"id", "created_at", "value"}, ", ")
}
