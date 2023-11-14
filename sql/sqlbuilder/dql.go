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
	Where      WhereClause
	Group      Group
	Having     HavingClause
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
	if l.Where.Valid() {
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
	if l.Where.Valid() {
		cs = append(cs, l.Where)
	}
	if l.Group != nil {
		cs = append(cs, l.Group)
	}
	if l.Having.Valid() {
		cs = append(cs, l.Having)
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

func buildConditions(name string, conditions []Condition, sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	if len(conditions) > 0 {
		err = WriteStringWithSpace(sqlWriter, name, level)
		if err != nil {
			return err
		}
		err = EndLine(sqlWriter, CompactLevel(level))
		if err != nil {
			return err
		}
		for i, c := range conditions {
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

type WhereClause struct {
	Conditions []Condition
}

func (c WhereClause) Valid() bool {
	return c.Conditions != nil
}

func (c WhereClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	return buildConditions("WHERE", c.Conditions, sqlWriter, argWriter, level)
}

type HavingClause struct {
	Conditions []Condition
}

func (c HavingClause) Valid() bool {
	return c.Conditions != nil
}

func (c HavingClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	return buildConditions("HAVING", c.Conditions, sqlWriter, argWriter, level)
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

type NamePosition int

const (
	NameFirst = 0 // name AS (SELECT ...)
	NameAfter = 1 // (SELECT ...) AS name
)

type Table struct {
	Clause       Clause
	Name         string
	NamePosition NamePosition
}

func (c Table) Valid() bool {
	return c.Clause != nil || c.Name != ""
}

func TableByName(name string) Table {
	return Table{Name: name}
}

func TableByClause(clause Clause) Table {
	return Table{Clause: clause}
}

func NameAsTable(name string, table Clause) Table {
	return Table{
		Clause:       table,
		Name:         name,
		NamePosition: NameFirst,
	}
}

func TableAsName(table Clause, name string) Table {
	return Table{
		Clause:       table,
		Name:         name,
		NamePosition: NameAfter,
	}
}

type From struct {
	Table Table
}

func (c From) Valid() bool {
	return c.Table.Valid()
}

func FromTable(clause Clause) From {
	return From{
		Table: TableByClause(clause),
	}
}

func FromTableName(name string) From {
	return From{
		Table: TableByName(name),
	}
}

func (c From) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	if c.Table.Clause != nil && c.Table.Name != "" {
		switch c.Table.NamePosition {
		case NameFirst:
			return ParseNameFirst(c.Table.Name, c.Table.Clause, sqlWriter, argWriter, level)
		case NameAfter:
			return ParseNameAfter(c.Table.Name, c.Table.Clause, sqlWriter, argWriter, level)
		default:
			panic("bad NamePosition")
		}
	}
	if c.Table.Clause != nil {
		return ParseSubTable(c.Table.Clause, sqlWriter, argWriter, level)
	}
	if c.Table.Name != "" {
		return ParseTableName(c.Table.Name, sqlWriter, argWriter, level)
	}
	return nil
}

func ParseTableName(name string, sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
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

func ParseSubTable(table Clause, sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
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

func ParseNameFirst(name string, clause Clause, sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	err = WriteStringWithSpace(sqlWriter, "FROM ", level)
	if err != nil {
		return err
	}
	err = WriteString(sqlWriter, name)
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
	err = clause.Parse(sqlWriter, argWriter, NextLevel(level))
	if err != nil {
		return err
	}
	err = WriteString(sqlWriter, ")")
	if err != nil {
		return err
	}
	err = EndLine(sqlWriter, CompactLevel(level))
	if err != nil {
		return err
	}
	return nil
}

func ParseNameAfter(name string, clause Clause, sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	err = WriteStringWithSpace(sqlWriter, "FROM (", level)
	if err != nil {
		return err
	}
	err = EndLine(sqlWriter, CompactLevel(level))
	if err != nil {
		return err
	}
	err = clause.Parse(sqlWriter, argWriter, NextLevel(level))
	if err != nil {
		return err
	}
	err = WriteString(sqlWriter, ") AS ")
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

type With interface {
	Clause
}

// WithClause generage WITH clause in SQL, its level arguments are always 0.
type WithClause struct {
	Tables []Table
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
		if t.Clause != nil {
			switch t.NamePosition {
			case NameFirst:
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
				err = t.Clause.Parse(sqlWriter, argWriter, NextLevel(level))
				if err != nil {
					return err
				}
				err = WriteString(sqlWriter, ")")
				if err != nil {
					return err
				}
			case NameAfter:
				err = WriteString(sqlWriter, "(")
				if err != nil {
					return err
				}
				err = EndLine(sqlWriter, CompactLevel(level))
				if err != nil {
					return err
				}
				err = t.Clause.Parse(sqlWriter, argWriter, 1)
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
			default:
				panic("bad NamePosition")
			}
		} else {
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
