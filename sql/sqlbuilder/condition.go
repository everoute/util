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
	OmitBrackets = true
	SaveBrackets = false
)

type AndCondition struct {
	L            Condition
	R            Condition
	OmitBrackets bool
}

func (c AndCondition) Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error {
	var err error
	if !c.OmitBrackets {
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
	if !c.OmitBrackets {
		err = WriteString(sqlWriter, ")")
		if err != nil {
			return err
		}
	}
	return nil
}

func And(l, r Condition, omitBrackets bool) AndCondition {
	return AndCondition{
		L:            l,
		R:            r,
		OmitBrackets: omitBrackets,
	}
}

type OrCondition struct {
	L            Condition
	R            Condition
	OmitBrackets bool
}

func (c OrCondition) Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error {
	var err error
	if !c.OmitBrackets {
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
	if !c.OmitBrackets {
		err = WriteString(sqlWriter, ")")
		if err != nil {
			return err
		}
	}
	return nil
}

func Or(l, r Condition, omitBrackets bool) OrCondition {
	return OrCondition{
		L:            l,
		R:            r,
		OmitBrackets: omitBrackets,
	}
}

type NotCondition struct {
	Condition    Condition
	OmitBrackets bool
}

func (c NotCondition) Parse(sqlWriter io.StringWriter, argWriter ArgWriter) error {
	var err error
	if !c.OmitBrackets {
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
	if !c.OmitBrackets {
		err = WriteString(sqlWriter, ")")
		if err != nil {
			return err
		}
	}
	return nil
}

func Not(condition Condition, omitBrackets bool) NotCondition {
	return NotCondition{
		Condition:    condition,
		OmitBrackets: omitBrackets,
	}
}
