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

const (
	Format  = 0
	Compact = -1
)

func NextLevel(level int) int {
	if level == -1 {
		return -1
	}
	return level + 1
}

func CompactLevel(level int) bool {
	return level < 0
}

type Clauses []Clause

func (cs *Clauses) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	for _, c := range *cs {
		if err := c.Parse(sqlWriter, argWriter, level); err != nil {
			return err
		}
	}
	return nil
}

const (
	AutoNewline = true
	DontNewline = false
)

func NewCustomClause(fn func(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error) CustomClause {
	return fn
}

type CustomClause func(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error

func (c CustomClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	return c(sqlWriter, argWriter, level)
}

func NewSimpleClause(autoEndline bool, SQL string, args ...Arg) Clause {
	return &SimpleClause{
		SQL:         SQL,
		Args:        args,
		AutoEndline: autoEndline,
	}
}

type SimpleClause struct {
	SQL         string
	Args        []Arg
	AutoEndline bool
}

func (c *SimpleClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	var err error
	err = WriteStringWithSpace(sqlWriter, c.SQL, level)
	if err != nil {
		return err
	}
	if c.AutoEndline {
		err = EndLine(sqlWriter, CompactLevel(level))
		if err != nil {
			return err
		}
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
	if CompactLevel(level) || CompactLevel(c.level) {
		return c.clause.Parse(sqlWriter, argWriter, Compact)
	}
	return c.clause.Parse(sqlWriter, argWriter, level+c.level)
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
