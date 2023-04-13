package builq_test

import (
	"fmt"
	"regexp"

	"github.com/cristalhq/builq"
)

func ExampleBuilder() {
	cols := builq.Columns{"foo", "bar"}

	var sb builq.Builder
	sb.Addf("SELECT %s FROM %s", cols, "table")
	sb.Addf("WHERE id = %$", 123)

	// this WILL NOT complile
	// var orClause = "OR id = %$"
	// sb.Addf(orClause, 42)

	// WILL compile
	const orClause2 = "OR id = %$"
	sb.Addf(orClause2, 42)

	query, args, err := sb.Build()
	panicIf(err)

	fmt.Printf("query:\n%v", query)
	fmt.Printf("args:\n%v", args)

	// Output:
	// query:
	// SELECT foo, bar FROM table
	// WHERE id = $1
	// OR id = $2
	// args:
	// [123 42]
}

func ExampleOnelineBuilder() {
	cols := builq.Columns{"foo", "bar"}

	var b builq.OnelineBuilder
	b.Addf("SELECT %s FROM %s", cols, "table")
	b.Addf("WHERE id = %$", 123)

	query, _, err := b.Build()
	panicIf(err)

	fmt.Print(query)

	// Output:
	// SELECT foo, bar FROM table WHERE id = $1
}

func ExampleBuilder_DebugBuild() {
	cols := builq.Columns{"foo", "bar"}

	var sb builq.Builder
	sb.Addf("SELECT %s FROM table", cols)
	sb.Addf("WHERE id = %$", 123)
	sb.Addf("OR id = %$ + %d", "42", 690)

	fmt.Printf("debug:\n%v", sb.DebugBuild())

	// Output:
	// debug:
	// SELECT foo, bar FROM table
	// WHERE id = 123
	// OR id = '42' + 690
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

func Example_query1() {
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

func Example_query2() {
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

func Example_query3() {
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

func Example_queryWhere() {
	filter := map[string]any{
		"name":     "the best",
		"category": []int{1, 2, 3},
		"pat":      regexp.MustCompile("pat+"),
		"prob":     0.42,
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
	if pat, ok := filter["pat"]; ok {
		b.Addf("AND page LIKE '%s'", pat)
	}
	if prob, ok := filter["prob"]; ok {
		b.Addf("AND prob < %s", prob)
	}
	if limit, ok := filter["limit"]; ok {
		b.Addf("LIMIT %d;", limit)
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
	// AND page LIKE 'pat+'
	// AND prob < 0.42
	// LIMIT 100;
	// args:
	// [the best 1 2 3]
}

func Example_slicePostgres() {
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

func Example_sliceMySQL() {
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

func Example_insertReturn() {
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

func Example_batchPostgres() {
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

func Example_batchMySQL() {
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

func Example_sliceInBatch() {
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
