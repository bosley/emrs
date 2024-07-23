package legate

type monotonic struct {
	value uint64
}

func (m *monotonic) get() uint64 {
	x := m.value
	m.value++
	return x
}

func (m *monotonic) GetExports() []Export {
	return []Export{
		// Only one function exported by this module
		Export{
			Name: "monotonic_get",
			Fn:   m.get,
		},
	}
}

// Creates a runtime object for legate that exposes
// a shated monotonic counter to all runtimes
// launched by legate (mostly for testing)
func NewModuleMonotonic() RuntimeModule {
	return &monotonic{
		value: 0,
	}
}
