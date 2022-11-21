package builq_test

import (
	"fmt"

	"github.com/cristalhq/builq"
)

func ExampleQuery1() {
	cols := builq.Columns{"foo, bar"}

	var b builq.Builder
	b.Addf("SELECT %s FROM %s", cols, "users").
		Addf("WHERE active IS TRUE").
		Addf("AND user_id = %$ OR user = %$", 42, "root")

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
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

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// SELECT foo, bar FROM users
	// WHERE
	// active = $1
	// AND user_id = $2
	// ORDER BY created_at
	// LIMIT 100;
	// args:
	// [true 42]
}

func ExampleQuery3() {
	var b builq.Builder
	b.Addf("SELECT * FROM foo").
		Addf("WHERE active IS TRUE").
		Addf("AND user_id = %$", 42).
		Addf("LIMIT 100;")

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// SELECT * FROM foo
	// WHERE active IS TRUE
	// AND user_id = $1
	// LIMIT 100;
	// args:
	// [42]
}

func ExampleQueryWhere() {
	filter := map[string]any{
		"name":     "the best",
		"category": []int{1, 2, 3},
		"page":     42,
		"limit":    100,
	}

	var b builq.Builder
	b.Addf("SELECT * FROM foo")
	b.Addf("WHERE active IS TRUE")

	if name, ok := filter["name"]; ok {
		b.Addf("AND name = %$", name)
	}
	if cat, ok := filter["category"]; ok {
		b.Addf("AND category IN (%+$)", cat)
	}
	if page, ok := filter["page"]; ok {
		b.Addf("AND page LIKE %$", page)
	}
	if limit, ok := filter["limit"]; ok {
		b.Addf("LIMIT %$;", limit)
	}

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// SELECT * FROM foo
	// WHERE active IS TRUE
	// AND name = $1
	// AND category IN ($2, $3, $4)
	// AND page LIKE $5
	// LIMIT $6;
	// args:
	// [the best 1 2 3 42 100]
}

func ExampleColumns() {
	columns := builq.Columns{"id", "created_at", "value"}
	params := []any{42, "right now", "just testing"}

	var b builq.Builder
	b.Addf("INSERT INTO %s (%s)", "table", columns)
	b.Addf("VALUES (%?, %?, %?);", params...)

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// INSERT INTO table (id, created_at, value)
	// VALUES (?, ?, ?);
	// args:
	// [42 right now just testing]
}

func ExampleSlicePostgres() {
	params := []any{42, true, "str"}

	var b builq.Builder
	b.Addf("INSERT INTO table (id, flag, name)")
	b.Addf("VALUES (%+$);", params)

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// INSERT INTO table (id, flag, name)
	// VALUES ($1, $2, $3);
	// args:
	// [42 true str]
}

func ExampleSliceMySQL() {
	params := []any{42, true, "str"}

	var b builq.Builder
	b.Addf("INSERT INTO table (id, flag, name)")
	b.Addf("VALUES (%+?);", params)

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// INSERT INTO table (id, flag, name)
	// VALUES (?, ?, ?);
	// args:
	// [42 true str]
}

func ExampleInsertReturn() {
	cols := builq.Columns{"id", "is_active", "name"}
	params := []any{true, "str"}

	var b builq.Builder
	b.Addf("INSERT INTO table (%s)", cols[1:]) // skip id column
	b.Addf("VALUES (%+$)", params)
	b.Addf("RETURNING %s;", cols)

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// INSERT INTO table (is_active, name)
	// VALUES ($1, $2)
	// RETURNING id, is_active, name;
	// args:
	// [true str]
}

func ExampleBatchPostgres() {
	params := [][]any{
		{42, true, "str"},
		{69, true, "noice"},
	}

	var b builq.Builder
	b.Addf("INSERT INTO table (id, flag, name)")
	b.Addf("VALUES %#$;", params)

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// INSERT INTO table (id, flag, name)
	// VALUES ($1, $2, $3), ($4, $5, $6);
	// args:
	// [42 true str 69 true noice]
}

func ExampleBatchMySQL() {
	params := [][]any{
		{42, true, "str"},
		{69, true, "noice"},
	}

	var b builq.Builder
	b.Addf("INSERT INTO table (id, flag, name)")
	b.Addf("VALUES %#?;", params)

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// INSERT INTO table (id, flag, name)
	// VALUES (?, ?, ?), (?, ?, ?);
	// args:
	// [42 true str 69 true noice]
}

func ExampleSliceInBatch() {
	params := [][]any{
		{42, []any{1, 2, 3}},
		{69, []any{4, 5, 6}},
	}

	var b builq.Builder
	b.Addf("INSERT INTO table (id, flag, name)")
	b.Addf("VALUES %#?;", params)

	query, args, err := b.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// INSERT INTO table (id, flag, name)
	// VALUES (?, ?), (?, ?);
	// args:
	// [42 [1 2 3] 69 [4 5 6]]
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
