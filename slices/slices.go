package slices

func ToAny[E any](s []E) []any {
	if s == nil {
		return nil
	}
	newSlice := make([]any, 0, len(s))
	for i := range s {
		newSlice = append(newSlice, s[i])
	}
	return newSlice
}

func ToInterface[E any](s []E) []interface{} {
	if s == nil {
		return nil
	}
	newSlice := make([]interface{}, 0, len(s))
	for i := range s {
		newSlice = append(newSlice, s[i])
	}
	return newSlice
}

func InterfaceToAny(s []interface{}) []any {
	return s
}

func AnyToInterface(s []any) []interface{} {
	return s
}

func FromAny[E any](s []any) []E {
	if s == nil {
		return nil
	}
	newSlice := make([]E, 0, len(s))
	for i := range s {
		newSlice = append(newSlice, s[i].(E))
	}
	return newSlice
}

func FromAnySafe[E any](s []any) ([]E, bool) {
	if s == nil {
		return nil, false
	}
	newSlice := make([]E, 0, len(s))
	for i := range s {
		e, ok := s[i].(E)
		if !ok {
			return nil, false
		}
		newSlice = append(newSlice, e)
	}
	return newSlice, true
}

func FromInterface[E any](s []interface{}) []E {
	if s == nil {
		return nil
	}
	newSlice := make([]E, 0, len(s))
	for i := range s {
		newSlice = append(newSlice, s[i].(E))
	}
	return newSlice
}

func FromInterfaceSafe[E any](s []interface{}) ([]E, bool) {
	if s == nil {
		return nil, false
	}
	newSlice := make([]E, 0, len(s))
	for i := range s {
		e, ok := s[i].(E)
		if !ok {
			return nil, false
		}
		newSlice = append(newSlice, e)
	}
	return newSlice, true
}
