package dataspec

// Sets the name of Spec to be used as the key of hcldec.ObjectSpec.
type KeyForObjectSpec struct {
	Spec
	key string
}

func (ns *KeyForObjectSpec) KeyForObjectSpec() string {
	return ns.key
}

// Gets the name of the item (block or attr).
// It may be different form the name this item will have in the Object.
func ItemName(s Spec) string {
	switch sT := s.getSpec().(type) {
	case *AttrSpec:
		return sT.Name
	case *BlockSpec:
		return sT.Name
	}
	return ""
}

// Specifies the key to be used in object spec
func UnderKey(key string, spec Spec) *KeyForObjectSpec {
	return &KeyForObjectSpec{
		key:  key,
		Spec: spec,
	}
}
