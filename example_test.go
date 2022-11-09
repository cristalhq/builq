package builq_test

import (
	"fmt"
	"time"

	"github.com/cristalhq/builq"
)

func ExampleQuery1() {
	cols := builq.Columns{"foo, bar"}

	var b builq.Builder
	b.Addf("SELECT %s FROM %s", cols, "users").
		Addf("WHERE active IS TRUE").
		Addf("AND user_id = %$ OR user = %$", 42, "root")

	query, args, err := b.Build()
	if err != nil {
		panic(err)
	}

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	//
	// query:
	// SELECT foo, bar FROM users
	// WHERE active IS TRUE
	// AND user_id = $1 OR user = $2
	// args:
	// [42 root]
}

func ExampleQuery2() {
	var b builq.Builder
	b.Addf("SELECT %s FROM %s", "foo, bar", "users")
	b.Addf("WHERE")
	b.Addf("active = %$", true)
	b.Addf("AND user_id = %$", 42)
	b.Addf("ORDER BY created_at")
	b.Addf("LIMIT 100;")

	query, _, err := b.Build()
	if err != nil {
		panic(err)
	}

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
	var b builq.Builder
	b.Addf("SELECT * FROM foo").
		Addf("WHERE active IS TRUE").
		Addf("AND user_id = %$", 42).
		Addf("LIMIT 100;")

	query, _, err := b.Build()
	if err != nil {
		panic(err)
	}

	fmt.Println(query)

	// Output:
	//
	// SELECT * FROM foo
	// WHERE active IS TRUE
	// AND user_id = $1
	// LIMIT 100;
}

func ExampleColumns() {
	columns := builq.Columns{"id", "created_at", "value"}
	args := []any{42, time.Now(), "just testing"}

	var b builq.Builder
	b.Addf("INSERT (%s) INTO %s", columns, "table")
	b.Addf("VALUES (%?, %?, %?);", args...) // TODO(junk1tm): should %a support slices?

	query, _, err := b.Build()
	if err != nil {
		panic(err)
	}

	fmt.Println(query)

	// Output:
	//
	// INSERT (id, created_at, value) INTO table
	// VALUES (?, ?, ?);
}

func ExampleSlice() {
	args := []any{42, true, "str"}

	var b builq.Builder
	b.Addf("INSERT (id, flag, name) INTO table")
	b.Addf("VALUES (%$);", args)
	query, _, err := b.Build()

	if err != nil {
		println(err.Error())
	}

	fmt.Println(query)

	// Output:
	//
	// INSERT (id, flag, name) INTO table
	// VALUES ($1);
}
