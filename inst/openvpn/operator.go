package openvpn

// Operator is a pipeline element structure.
type Operator struct {
	Name   string
	Run    func(*OpenVPN) error
	Cancel func(*OpenVPN) error
}

// NewOperator creates a new Operator instance.
func NewOperator(name string, run func(*OpenVPN) error,
	cancel func(*OpenVPN) error) *Operator {
	if cancel == nil {
		cancel = func(in *OpenVPN) error { return nil }
	}
	return &Operator{Name: name, Run: run, Cancel: cancel}
}
