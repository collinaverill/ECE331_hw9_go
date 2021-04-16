Plot some COVID data from an SQLite3 database
Uses https://github.com/mattn/go-sqlite3 and CGO

be sure to build with `CGO_ENABLED=1`

eg:
`CGO_ENABLED=1 go run main.go`

See [simple.go](https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go) for more examples