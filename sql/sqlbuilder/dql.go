package sqlbuilder

import (
	"io"
)

// The template of Data Query Language
// It is also possible to implement DQL through Clauses
type DQL struct {
	With       With
	Select     Select
	From       From
	Where      []Condition
	Group      Group
	Order      Order
	Limit      Limit
	Additional Clauses
}

func (l *DQL) Clauses() Clauses {
	count := 1 // The SELECT clause is MUST.
	if l.With != nil {
		count++
	}
	if l.From.Valid() {
		count++
	}
	if l.Where != nil {
		count++
	}
	if l.Group != nil {
		count++
	}
	if l.Order != nil {
		count++
	}
	if l.Limit != nil {
		count++
	}
	count += len(l.Additional)
	cs := make([]Clause, 0, count)
	if l.With != nil {
		cs = append(cs, l.With)
	}
	cs = append(cs, &l.Select)
	if l.From.Valid() {
		cs = append(cs, l.From)
	}
	if l.Where != nil {
		cs = append(cs, &whereClause{Where: l.Where})
	}
	if l.Group != nil {
		cs = append(cs, l.Group)
	}
	if l.Order != nil {
		cs = append(cs, l.Order)
	}
	if l.Limit != nil {
		cs = append(cs, l.Limit)
	}
	if l.Additional != nil {
		cs = append(cs, l.Additional...)
	}
	return cs
}

type whereClause struct {
	Where []Condition
}

func (c *whereClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	if len(c.Where) > 0 {
		err = WriteStringWithSpace(sqlWriter, "WHERE", level)
		if err != nil {
			return err
		}
		err = EndLine(sqlWriter, CompactLevel(level))
		if err != nil {
			return err
		}
		for i, c := range c.Where {
			if i != 0 {
				err = WriteStringWithSpace(sqlWriter, "AND ", NextLevel(level))
				if err != nil {
					return err
				}
				err = c.Parse(sqlWriter, argWriter)
				if err != nil {
					return err
				}
			} else {
				err = WriteSpace(sqlWriter, NextLevel(level))
				if err != nil {
					return err
				}
				err = c.Parse(sqlWriter, argWriter)
				if err != nil {
					return err
				}
			}
			err = EndLine(sqlWriter, CompactLevel(level))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *DQL) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	cs := l.Clauses()
	return cs.Parse(sqlWriter, argWriter, level)
}

type Select struct {
	Columns []string
	// Args are valid if the Columns is not a empty slice or nil.
	Args []Arg
}

func (c *Select) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	if len(c.Columns) == 0 {
		err = WriteStringWithSpace(sqlWriter, "SELECT *", level)
		if err != nil {
			return err
		}
		err = EndLine(sqlWriter, CompactLevel(level))
		if err != nil {
			return err
		}
	} else {
		err = WriteStringWithSpace(sqlWriter, "SELECT", level)
		if err != nil {
			return err
		}
		err = EndLine(sqlWriter, CompactLevel(level))
		if err != nil {
			return err
		}
		for i, col := range c.Columns {
			err = WriteStringWithSpace(sqlWriter, col, NextLevel(level))
			if err != nil {
				return err
			}
			if i != len(c.Columns)-1 {
				err = WriteString(sqlWriter, ",")
				if err != nil {
					return err
				}
			}
			err = EndLine(sqlWriter, CompactLevel(level))
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

type From struct {
	Table Clause
	Name  string
}

func (c From) Valid() bool {
	return c.Name != "" || c.Table != nil
}

func FromName(name string) From {
	return From{Name: name}
}

func FromTable(table Clause) From {
	return From{Table: table}
}

func (c From) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	if c.Name != "" {
		return parseFromTableName(c.Name, sqlWriter, argWriter, level)
	}
	if c.Table != nil {
		return parseFromSubTable(c.Table, sqlWriter, argWriter, level)
	}
	return nil
}

func parseFromTableName(name string, sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	err = WriteStringWithSpace(sqlWriter, "FROM ", level)
	if err != nil {
		return err
	}
	err = WriteString(sqlWriter, name)
	if err != nil {
		return err
	}
	err = EndLine(sqlWriter, CompactLevel(level))
	if err != nil {
		return err
	}
	return nil
}

func parseFromSubTable(table Clause, sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	err = WriteStringWithSpace(sqlWriter, "FROM (", level)
	if err != nil {
		return err
	}
	err = EndLine(sqlWriter, CompactLevel(level))
	if err != nil {
		return err
	}
	err = table.Parse(sqlWriter, argWriter, NextLevel(level))
	if err != nil {
		return err
	}
	err = WriteStringWithSpace(sqlWriter, ")", level)
	if err != nil {
		return err
	}
	err = EndLine(sqlWriter, CompactLevel(level))
	if err != nil {
		return err
	}
	return nil
}

type NamePosition int

const (
	NameFirst = 0 // name AS (SELECT ...)
	NameAfter = 1 // (SELECT ...) AS name
)

type NamedTable struct {
	Name  string
	Table Clause
}

type With interface {
	Clause
}

// WithClause generage WITH clause in SQL, its level arguments are always 0.
type WithClause struct {
	Tables       []NamedTable
	NamePosition NamePosition
}

func (c *WithClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	if !CompactLevel(level) {
		level = 0
	}
	var err error
	err = WriteString(sqlWriter, "WITH")
	if err != nil {
		return err
	}
	err = EndLine(sqlWriter, CompactLevel(level))
	if err != nil {
		return err
	}
	for i, t := range c.Tables {
		if c.NamePosition == NameFirst {
			err = WriteString(sqlWriter, t.Name)
			if err != nil {
				return err
			}
			err = WriteString(sqlWriter, " AS (")
			if err != nil {
				return err
			}
			err = EndLine(sqlWriter, CompactLevel(level))
			if err != nil {
				return err
			}
			err = t.Table.Parse(sqlWriter, argWriter, NextLevel(level))
			if err != nil {
				return err
			}
			err = WriteString(sqlWriter, ")")
			if err != nil {
				return err
			}
		} else {
			err = WriteString(sqlWriter, "(")
			if err != nil {
				return err
			}
			err = EndLine(sqlWriter, CompactLevel(level))
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
		err = EndLine(sqlWriter, CompactLevel(level))
		if err != nil {
			return err
		}
	}
	return nil
}

type Group interface {
	Clause
}

func MakeGroupby(value string, args ...Arg) Group {
	return NewSimpleClause(AutoNewline, "GROUP BY "+value, args...)
}

type Order interface {
	Clause
}

func MakeOrderby(value string, args ...Arg) Order {
	return NewSimpleClause(AutoNewline, "ORDER BY "+value, args...)
}

type Limit interface {
	Clause
}

func MakeLimit(value string, args ...Arg) Order {
	return NewSimpleClause(AutoNewline, "LIMIT "+value, args...)
}
