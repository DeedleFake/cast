package usage

// Interface represents a network interface present on the system.
type Interface struct {
	inter
}

// OpenInterface opens the interface with the specified name.
func OpenInterface(name string) (*Interface, error) {
	return openInterface(name)
}

// In returns the number of bytes received by the interface since the
// last time that it was called, or since the interface was opened if
// it hasn't been called before.
func (i *Interface) In() (uint64, error) {
	return i.inter.In()
}

// Out returns the number of bytes sent by the interface since the
// last time that it was called, or since the interface was opened if
// it hasn't been called before.
func (i *Interface) Out() (uint64, error) {
	return i.inter.Out()
}
