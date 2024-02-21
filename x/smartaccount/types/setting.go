package types

import fmt "fmt"

// NewParams creates a new Params object
func NewSetting(ownerAddr string) Setting {
	return Setting{
		Owner:    ownerAddr,
		Fallback: true,
	}
}

func NewSettings() []*Setting {
	return make([]*Setting, 0)
}

func DefaultSettings() []*Setting {
	return NewSettings()
}

func (s Setting) Validate() error {
	if s.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}

	if (s.Authorization == nil || len(s.Authorization) == 0) && !s.Fallback {
		return fmt.Errorf("fallback must be true without any authorization")
	}

	return nil
}
