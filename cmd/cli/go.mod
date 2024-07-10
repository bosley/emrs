module github.com/bosley/cmd/cli

go 1.22.3

replace github.com/bosley/emrs/badger => ../../badger

replace github.com/bosley/emrs/app => ../../app

require (
	github.com/bosley/emrs/app v0.0.0-00010101000000-000000000000
	github.com/bosley/emrs/badger v0.0.0-00010101000000-000000000000
)

require golang.org/x/crypto v0.24.0 // indirect
