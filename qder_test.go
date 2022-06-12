package qder_test

import (
	"database/sql"
	"fmt"

	"github.com/cristalhq/qder"
)

func ExampleQuery1() {
	q := qder.Newf("SELECT %s FROM %s", "foo, bar", "users")
	q.Append("WHERE")
	q.Append("active IS TRUE")
	x := q.AddParam(42)
	q.Append("AND user_id =", x, "OR invited_by =", x)

	doQuery(q.Query(), q.Args()...)

	// Output:
	//
	// SELECT foo, bar FROM users
	// WHERE
	// active IS TRUE
	// AND user_id = $1 OR invited_by = $1
}

func ExampleQuery2() {
	q := qder.Newf("SELECT %s FROM %s", "foo, bar", "users")
	q.Append("WHERE")
	q.Add("active = ", true)
	q.Add("AND user_id = ", 42)
	q.Append("ORDER BY created_at")
	q.Append("LIMIT 100;")

	doQuery(q.Query(), q.Args()...)

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
	q := qder.Newf("SELECT * FROM foo")
	q.Append("WHERE active IS TRUE")
	q.Append("AND user_id = " + q.AddParam(42))
	q.Append("LIMIT 100;")

	doQuery(q.Query(), q.Args()...)

	// Output:
	//
	// SELECT * FROM foo
	// WHERE active IS TRUE
	// AND user_id = $1
	// LIMIT 100;
}

func doQuery(query string, args ...interface{}) {
	fmt.Println(query)
	if false { // because we don't have db to test it
		var db *sql.DB
		db.Query(query, args...)
	}
}
