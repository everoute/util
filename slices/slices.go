package slices

import (
	std "slices"
)

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

func InterfaceAsAny(s []interface{}) []any {
	return s
}

func AnyAsInterface(s []any) []interface{} {
	return s
}

func InterfaceToAny(s []interface{}) []any {
	return std.Clone(([]any)(s))
}

func AnyToInterface(s []any) []interface{} {
	return std.Clone(([]interface{})(s))
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

func Concat[E any](s1, s2 []E) []E {
	if s1 == nil {
		return std.Clone(s2)
	}
	if s2 == nil {
		return std.Clone(s1)
	}
	newSlice := make([]E, 0, len(s1)+len(s2))
	newSlice = append(newSlice, s1...)
	newSlice = append(newSlice, s2...)
	return newSlice
}

func AppendAny[E any](s1 []E, s2 []any) []E {
	for i := range s2 {
		s1 = append(s1, s2[i].(E))
	}
	return s1
}

func AppendAnySafe[E any](s1 []E, s2 []any) ([]E, bool) {
	for i := range s2 {
		e, ok := s2[i].(E)
		if !ok {
			return nil, false
		}
		s1 = append(s1, e)
	}
	return s1, true
}

func AppendToAny[E any](s1 []any, s2 []E) []any {
	for i := range s2 {
		s1 = append(s1, s2[i])
	}
	return s1
}

func First[E any](s []E) (r E) {
	if len(s) == 0 {
		return
	}
	return s[0]
}

func Last[E any](s []E) (r E) {
	l := len(s)
	if l == 0 {
		return
	}
	return s[l-1]
}

func Map[S, D any](s []S, fn func(S) D) []D {
	l := len(s)
	newSlice := make([]D, 0, l)
	for i := 0; i < l; i++ {
		newSlice = append(newSlice, fn(s[i]))
	}
	return newSlice
}

func MapWithIdx[S, D any](s []S, fn func(S, int) D) []D {
	l := len(s)
	newSlice := make([]D, 0, l)
	for i := 0; i < l; i++ {
		newSlice = append(newSlice, fn(s[i], i))
	}
	return newSlice
}
