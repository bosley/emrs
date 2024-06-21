module github.com/bosley/emrs

go 1.22.3

replace emrs/core => ./core

replace emrs/webui => ./webui

require emrs/core v0.0.0-00010101000000-000000000000

require gopkg.in/yaml.v3 v3.0.1

require emrs/webui v0.0.0-00010101000000-000000000000 // indirect
