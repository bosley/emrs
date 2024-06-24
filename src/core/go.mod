module emrs/core

go 1.22.3

replace emrs/datastore => ../datastore

require emrs/datastore v0.0.0-00010101000000-000000000000

require github.com/mattn/go-sqlite3 v1.14.22 // indirect
