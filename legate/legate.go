package legate

import (
	"time"
)

const (
	defaultTargetName = "legate.yaml"
)

const (
	MemPage1Gb = 16384
	MemPage2Gb = MemPage1Gb * 2
	MemPage3Gb = MemPage1Gb * 3
	MemPage4Gb = MemPage1Gb * 4

	MemPage500MB = MemPage1Gb / 2
	MemPage255MB = MemPage1Gb / 4

	MemPageMax = MemPage4Gb
)

type Opts struct {
	Ttl     time.Duration `yaml:ttl`   //  Max time target is permitted to process before term
	Pages   uint32        `yaml:pages` //  Max memory pages at 2^16 bytes each (65.536 KB)
	Targets []string      `targets`    //  Paths to target configurations
}

type Target struct {
	Name    string `yaml:name`
	Kind    string `yaml:kind`    // tinygo, zig, rust, .. (later) c, c++ etc
	Handler string `yaml:handler` // function name that takes array of bytes
	// There will be more target-specific config so we'll keep a struct
}
