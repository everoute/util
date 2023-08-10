package sqlbuilder

import (
	"fmt"
	"io"
)

/*
The clause in sql
level: the level
*/
type Clause interface {
	Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error
}

type Clauses []Clause

func (cs *Clauses) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	for _, c := range *cs {
		if err := c.Parse(sqlWriter, argWriter, level); err != nil {
			return err
		}
		if err := EndLine(sqlWriter); err != nil {
			return err
		}
	}
	return nil
}

func NewSimpleClause(SQL string, args ...Arg) Clause {
	return &SimpleClause{
		SQL:  SQL,
		Args: args,
	}
}

type SimpleClause struct {
	SQL  string
	Args []Arg
}

func (c *SimpleClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	err = WriteSpace(sqlWriter, level)
	if err != nil {
		return err
	}
	err = WriteString(sqlWriter, c.SQL)
	if err != nil {
		return err
	}
	err = WriteArgs(argWriter, c.Args...)
	if err != nil {
		return err
	}
	return nil
}

func AddClauseLevel(c Clause, level int) Clause {
	return &addLeveledClause{clause: c, level: level}
}

type addLeveledClause struct {
	clause Clause
	level  int
}

func (c *addLeveledClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	if err := WriteSpace(sqlWriter, 1); err != nil {
		return err
	}
	return c.clause.Parse(sqlWriter, argWriter, level)
}

func WriteString(writer io.StringWriter, str string) error {
	if n, err := writer.WriteString(str); err != nil {
		return err
	} else if n != len(str) {
		return fmt.Errorf("write string failed, expect: %d, real: %d", len(str), n)
	}
	return nil
}

func WriteStringWithSpace(writer io.StringWriter, str string, level int) error {
	if err := WriteSpace(writer, level); err != nil {
		return err
	}
	if n, err := writer.WriteString(str); err != nil {
		return err
	} else if n != len(str) {
		return fmt.Errorf("write string failed, expect: %d, real: %d", len(str), n)
	}
	return nil
}
