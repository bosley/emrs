module github.com/bosley/emrs

go 1.22.3

require (
	emrs/badger v0.0.0-00010101000000-000000000000
	emrs/core v0.0.0-00010101000000-000000000000
)

require golang.org/x/crypto v0.24.0 // indirect

replace emrs/badger => ./badger

replace emrs/core => ./core
