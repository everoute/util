package sqlbuilder

import (
	"io"
)

// Data query language
type DQL struct {
	With   With
	Select Select
	From   From
	Where  []Condition
	Group  Group
	Order  Order
	Limit  Limit
}

func (l *DQL) Clauses() Clauses {
	cs := make([]Clauses, 0)
	cs = append(cs)
	return nil
}

type whereClause struct {
	Where []Condition
}

func (c *whereClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	if len(c.Where) > 0 {
		err = WriteStringWithSpace(sqlWriter, "WHERE\n", level)
		if err != nil {
			return err
		}
		for i, c := range c.Where {
			if i != 0 {
				err = WriteStringWithSpace(sqlWriter, "AND ", level+1)
				if err != nil {
					return err
				}
			}
			err = c.Parse(sqlWriter, argWriter, 0)
			if err != nil {
				return err
			}
			err = EndLine(sqlWriter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *DQL) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	err = l.With.Parse(sqlWriter, argWriter, level)
	if err != nil {
		return err
	}
	// Write SELECT clause
	err = l.Select.Parse(sqlWriter, argWriter, level)
	if err != nil {
		return err
	}
	// Write From clause
	err = l.From.Parse(sqlWriter, argWriter, level)
	if err != nil {
		return err
	}
	// Write Where clause
	var whereC = whereClause{
		Where: l.Where,
	}
	err = whereC.Parse(sqlWriter, argWriter, level)
	if err != nil {
		return err
	}
	// Write Group by clause
	if l.Group != nil {
		l.Group.Parse(sqlWriter, argWriter, level)
	}
	// Write Orger by clause
	if l.Order != nil {
		l.Order.Parse(sqlWriter, argWriter, level)
	}
	// Write Limit clause
	if l.Limit != nil {
		l.Limit.Parse(sqlWriter, argWriter, level)
	}
	return nil
	// return errors.New("not implemented")
}

type Select struct {
	Columns []string
	// Args are valid if the Columns is not a empty slice or nil.
	Args []Arg
}

func (c *Select) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	if len(c.Columns) == 0 {
		err = WriteStringWithSpace(sqlWriter, "SELECT *\n", level)
		if err != nil {
			return err
		}
	} else {
		err = WriteStringWithSpace(sqlWriter, "SELECT\n", level)
		if err != nil {
			return err
		}
		for i, col := range c.Columns {
			err = WriteStringWithSpace(sqlWriter, col, level+1)
			if err != nil {
				return err
			}
			if i != len(c.Columns)-1 {
				err = WriteString(sqlWriter, ",")
				if err != nil {
					return err
				}
			}
			err = EndLine(sqlWriter)
			if err != nil {
				return err
			}
		}
		err = WriteArgs(argWriter, c.Args...)
		if err != nil {
			return err
		}
	}
	return nil
}

type Table interface {
	Clause
}

type From struct {
	Table Table
	Name  string
}

func (c *From) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	if c.Table != nil {
		err = WriteStringWithSpace(sqlWriter, "FROM\n", level)
		if err != nil {
			return err
		}
		err = c.Table.Parse(sqlWriter, argWriter, level+1)
		if err != nil {
			return err
		}

	} else {
		err = WriteStringWithSpace(sqlWriter, "FROM ", level)
		if err != nil {
			return err
		}
		err = WriteString(sqlWriter, c.Name)
		if err != nil {
			return err
		}
		err = EndLine(sqlWriter)
		if err != nil {
			return err
		}
	}

	return nil
}

type NamePosition int

const (
	NameFirst = 0
	NameAfter = 1
)

type NamedTable struct {
	Name  string
	Table Table
}

type Condition interface {
	Clause
}

type With interface {
	Clause
}

// WithClause generage WITH clause in SQL, its level arguments are always 0.
type WithClause struct {
	Tables       []NamedTable
	NamePosition NamePosition
}

func (c *WithClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, _ int) error {
	var err error
	err = WriteString(sqlWriter, "WITH\n")
	if err != nil {
		return err
	}
	for i, t := range c.Tables {
		if c.NamePosition == NameFirst {
			err = WriteString(sqlWriter, t.Name)
			if err != nil {
				return err
			}
			err = WriteString(sqlWriter, " AS (\n")
			if err != nil {
				return err
			}
			err = t.Table.Parse(sqlWriter, argWriter, 1)
			if err != nil {
				return err
			}
			err = WriteString(sqlWriter, ")")
			if err != nil {
				return err
			}
		} else {
			err = WriteString(sqlWriter, "(\n")
			if err != nil {
				return err
			}
			err = t.Table.Parse(sqlWriter, argWriter, 1)
			if err != nil {
				return err
			}
			err = WriteString(sqlWriter, ") AS ")
			if err != nil {
				return err
			}
			err = WriteString(sqlWriter, t.Name)
			if err != nil {
				return err
			}
		}
		if i != len(c.Tables)-1 {
			err = WriteString(sqlWriter, ",")
			if err != nil {
				return err
			}
		}
		err = EndLine(sqlWriter)
		if err != nil {
			return err
		}
	}
	return nil
}

type Group interface {
	Clause
}

func MakeGroupby(customize bool, value string, args ...Arg) Group {
	if customize {
		return NewSimpleClause(value+"\n", args...)
	} else {
		return NewSimpleClause("GROUP BY "+value+"\n", args...)
	}
}

type Order interface {
	Clause
}

func MakeOrderby(customize bool, value string, args ...Arg) Order {
	if customize {
		return NewSimpleClause(value+"\n", args...)
	} else {
		return NewSimpleClause("ORDER BY "+value+"\n", args...)
	}
}

type Limit interface {
	Clause
}

func MakeLimit(customize bool, value string, args ...Arg) Order {
	if customize {
		return NewSimpleClause(value+"\n", args...)
	} else {
		return NewSimpleClause("LIMIT "+value+"\n", args...)
	}
}
