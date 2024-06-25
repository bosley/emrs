module emrs/core

go 1.22.3

replace emrs/datastore => ../datastore
replace emrs/badger => ../badger

require emrs/datastore v0.0.0-00010101000000-000000000000
