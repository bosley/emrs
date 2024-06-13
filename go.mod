module github.com/bosley/emrs

go 1.22.3

replace internal/reaper => ./internal/reaper

replace internal/vault => ./internal/vault

require (
	github.com/bosley/nerv-go v0.1.1
	internal/reaper v0.0.0-00010101000000-000000000000
	internal/vault v0.0.0-00010101000000-000000000000
)
