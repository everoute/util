package sqlbuilder

import "io"

type Condition interface {
	Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error
}

// Use SimpleCondition please if a simple condition is need.
// Such as IN, LIKE, BETWEEN, =, >, <, >=, <=, <> etc..
type SimpleCondition struct {
	Str  string
	Args []Arg
}

func (c SimpleCondition) Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error {
	var err error
	err = WriteString(sqlWriter, c.Str)
	if err != nil {
		return err
	}
	err = WriteArgs(argWriter, c.Args...)
	if err != nil {
		return err
	}
	return nil
}

func NewCondition(str string, args ...Arg) SimpleCondition {
	return SimpleCondition{
		Str:  str,
		Args: args,
	}
}

type CustomCondition func(sqlWriter io.StringWriter, argWriter ArgWriter) error

func NewCustomCondition(fn func(sqlWriter io.StringWriter, argWriter ArgWriter) error) CustomCondition {
	return fn
}

func (c CustomCondition) Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error {
	return c(sqlWriter, argWriter)
}

const (
	SaveBrackets = true
	OmitBrackets = false
)

type BracketedCondition struct {
	Condition Condition
}

func (c BracketedCondition) Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error {
	var err error
	err = WriteString(sqlWriter, "(")
	if err != nil {
		return err
	}
	err = c.Condition.Parse(sqlWriter, argWriter)
	if err != nil {
		return err
	}
	err = WriteString(sqlWriter, ")")
	if err != nil {
		return err
	}
	return nil
}

func Bracket(condition Condition) Condition {
	return BracketedCondition{Condition: condition}
}

func BracketIf(condition Condition, bracket bool) Condition {
	if bracket {
		return BracketedCondition{Condition: condition}
	}
	return condition
}

type AndCondition struct {
	L       Condition
	R       Condition
	Bracket bool
}

func (c AndCondition) Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error {
	var err error
	if c.Bracket {
		err = WriteString(sqlWriter, "(")
		if err != nil {
			return err
		}
	}
	err = c.L.Parse(sqlWriter, argWriter)
	if err != nil {
		return err
	}
	err = WriteString(sqlWriter, " AND ")
	if err != nil {
		return err
	}
	err = c.R.Parse(sqlWriter, argWriter)
	if err != nil {
		return err
	}
	if c.Bracket {
		err = WriteString(sqlWriter, ")")
		if err != nil {
			return err
		}
	}
	return nil
}

func And(l, r Condition, bracket bool) Condition {
	switch {
	case l != nil && r != nil:
		return AndCondition{
			L:       l,
			R:       r,
			Bracket: bracket,
		}
	case l != nil && r == nil:
		return BracketIf(l, bracket)
	case l == nil && r != nil:
		return BracketIf(r, bracket)
	default:
		return nil
	}
}

type OrCondition struct {
	L       Condition
	R       Condition
	Bracket bool
}

func (c OrCondition) Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error {
	var err error
	if c.Bracket {
		err = WriteString(sqlWriter, "(")
		if err != nil {
			return err
		}
	}
	err = c.L.Parse(sqlWriter, argWriter)
	if err != nil {
		return err
	}
	err = WriteString(sqlWriter, " OR ")
	if err != nil {
		return err
	}
	err = c.R.Parse(sqlWriter, argWriter)
	if err != nil {
		return err
	}
	if c.Bracket {
		err = WriteString(sqlWriter, ")")
		if err != nil {
			return err
		}
	}
	return nil
}

func Or(l, r Condition, bracket bool) Condition {
	switch {
	case l != nil && r != nil:
		return OrCondition{
			L:       l,
			R:       r,
			Bracket: bracket,
		}
	case l != nil && r == nil:
		return BracketIf(l, bracket)
	case l == nil && r != nil:
		return BracketIf(r, bracket)
	default:
		return nil
	}
}

type NotCondition struct {
	Condition Condition
	Bracket   bool
}

func (c NotCondition) Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error {
	var err error
	if c.Bracket {
		err = WriteString(sqlWriter, "(")
		if err != nil {
			return err
		}
	}
	err = WriteString(sqlWriter, "NOT ")
	if err != nil {
		return err
	}
	err = c.Condition.Parse(sqlWriter, argWriter)
	if err != nil {
		return err
	}
	if c.Bracket {
		err = WriteString(sqlWriter, ")")
		if err != nil {
			return err
		}
	}
	return nil
}

func Not(condition Condition, bracket bool) Condition {
	return NotCondition{
		Condition: condition,
		Bracket:   bracket,
	}
}
