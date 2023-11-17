# sql2go

<strong>sql2go</strong> is a lightweight [Go](https://go.dev/) library designed to parse SQL files into corresponding Golang structs.

## Installation

```bash
go get -u github.com/phenpessoa/sql2go
```

## Usage

### API

The library exposes a single function:

```go
func Parse[T any](dst *T, r io.Reader) error
```

- `dst`: A struct with string fields representing query names in the SQL file.
- `r`: An io.Reader pointing to the SQL file to be parsed.

### Syntax

The expected syntax in the SQL file is as follows:

```sql
-- name: QueryName
SELECT * FROM table_name;
```

This lib does not validate the query. It just simply passes the query as is (minus whitespaces) to the according field in the struct.

Check the [evaluator_test](https://github.com/phenpessoa/sql2go/blob/main/internal/parser/evaluator_test.go) file and the files in the [files](https://github.com/phenpessoa/sql2go/tree/main/internal/testdata/files) dir, to better understand how the file will be parsed.

On the tests there are a lot of edge cases that you can check.

### Example

queries.sql:

```sql
-- name: GetUsers
SELECT * FROM users;

-- name: DeleteUsers
-- Warning! This will erase all users in the system
DELETE FROM users;

-- name: DeleteUser
-- Deletes a specific user from the system
DELETE FROM users
    -- Just to demonstrate whitespace and in-query comments
    WHERE id = 123;
```

main.go:

```go
package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/phenpessoa/sql2go"
)

//go:embed queries.sql
var queryFile string

type Queries struct {
	GetUsers    string
	DeleteUsers string
	DeleteUser  string
}

var queries Queries

func main() {
	if err := sql2go.Parse(&queries, strings.NewReader(queryFile)); err != nil {
		panic(err)
	}

	// Now, queries will be populated with SQL queries.

	fmt.Println(
		queries.GetUsers == "SELECT * FROM users;",
	) // true

	fmt.Println(
		queries.DeleteUsers == "DELETE FROM users;",
	) // true

	fmt.Println(
		queries.DeleteUser == "DELETE FROM users\nWHERE id = 123;",
	) // true
}
```

## Why not sqlc?

[sqlc](https://github.com/sqlc-dev/sqlc) is amazing. But it comes with 2 drawbacks, imo.

- It it not 100% compatible with windows
  - this can be solved with Docker
- It doesn't work flawlessly for complex queries and schemas
  - The fact that it checks all the queries, and it has to parse the schema, while great, is a problem when your schema/query is complex. For example: it does not handle enums in postgresql properly and will parse the type as an `interface{}`.

On the other hand, it writes a lot of boilerplate for you, which is amazing. With sql2go you will still have to write the database code manually. The main goal is to keep syntax highlighting by using a .sql file and to prevent the headaches when having to use a backtick on your query. Go's strings are to blame here. There is [a proposal](https://github.com/golang/go/issues/32590) to improve the experience, but it has not been accepted yet.
