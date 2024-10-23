package validate

func lookupRefType(ref string, v *Validator) string {
	if t := v.Select().ID(ref).First(); t != nil {
		return t.object.Type()
	}

	return ""
}
