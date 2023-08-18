package sqlbuilder

import (
	"fmt"
	"io"
	"math"
	"strings"
)

type Arg any

type ArgWriter interface {
	WriteArg(arg Arg) error
}

func WriteArgs(writer ArgWriter, args ...Arg) error {
	for _, arg := range args {
		if err := writer.WriteArg(arg); err != nil {
			return err
		}
	}
	return nil
}

/*
The clause in sql
sqlWriter: A interface implemented WriteString, the SQL statements will be wrote into it.
argWriter: A interface implemented WriteArg, the arguments will be written in the same order as the SQL.
level: The level describes the indent level at the beginning of each line, the number of indented Spaces is level*2.
Note: The tool does not check whether the parameters match the parameters of the SQL, because different databases support different formats， please check it manually.
*/
type Clause interface {
	Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error
}

const (
	// The default value of the level parameter when parsing the SQL.
	Format  = 0           // Format the SQL
	Compact = math.MinInt // Do not format the SQL
)

// If you need to increase indentation, get the next Level value with it.
func NextLevel(level int) int {
	if CompactLevel(level) {
		return Compact
	}
	return level + 1
}

// Determine whether the level value in compact mode
func CompactLevel(level int) bool {
	return level < 0
}

// The substitute of []Clause, it implemented a Parse method
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

type CustomClause func(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error

func NewCustomClause(fn func(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error) CustomClause {
	return fn
}

func (c CustomClause) Parse(sqlWriter io.StringWriter, argWriter ArgWriter, level int) error {
	return c(sqlWriter, argWriter, level)
}

func NewSimpleClause(autoEndline bool, sql string, args ...Arg) Clause {
	return &SimpleClause{
		SQL:         sql,
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

const (
	EOL   = "\n"
	Space = " "
)

// End the line with EOL or Space
func EndLine(sqlWriter io.StringWriter, withSpace bool) error {
	if withSpace {
		if n, err := sqlWriter.WriteString(Space); err != nil {
			return err
		} else if n != 1 {
			return fmt.Errorf("write endline failed")
		}
		return nil
	}
	if n, err := sqlWriter.WriteString(EOL); err != nil {
		return err
	} else if n != 1 {
		return fmt.Errorf("write endline failed")
	}
	return nil
}

// Write indentation with space.
// If you want to add indentation when writing a string, then WriteStringWithSpace is a better choice.
func WriteSpace(sqlWriter io.StringWriter, level int) error {
	// Redundant judgments will be optimized in GetSpace.
	if CompactLevel(level) {
		return nil
	}
	space := GetSpace(level)
	if n, err := sqlWriter.WriteString(space); err != nil {
		return err
	} else if n != len(space) {
		return fmt.Errorf("write space(level:%d) failed", level)
	}
	return nil
}

const (
	singleSpace       = "  "
	lenSingleSpace    = len(singleSpace)
	optimizeSpaces    = "                                                              "
	lenOptimizeSpaces = len(optimizeSpaces)
	maxOptimizeLevel  = lenOptimizeSpaces / lenSingleSpace
)

// Get the Spaces corresponding to the level.
// And custom the indentation strategy is not supported currently.
// When you want to write indentation, you should always call WriteSpace instead of GetSpace.
func GetSpace(level int) string {
	if CompactLevel(level) {
		return ""
	}
	// Optimize the vast majority of scenarios
	if level >= 0 && level <= maxOptimizeLevel {
		return optimizeSpaces[:level*lenSingleSpace]
	}
	return strings.Repeat(singleSpace, level)
}
