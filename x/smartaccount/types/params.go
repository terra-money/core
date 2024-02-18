package types

// NewParams creates a new Params object
func NewParams() Params {
	return Params{}
}

func DefaultParams() Params {
	return Params{}
}

func (p Params) Validate() error {
	return nil
}
