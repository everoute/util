package sqlbuilder_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/everoute/util/sql/sqlbuilder"
)

type fixedBuilder struct {
	builder strings.Builder
	size    int
	cap     int
}

func newFixedBuilder(cap int) *fixedBuilder {
	return &fixedBuilder{
		builder: strings.Builder{},
		size:    0,
		cap:     cap,
	}
}

type ArgWriter struct {
	Args []sqlbuilder.Arg
}

func NewArgWriter(cap int) *ArgWriter {
	return &ArgWriter{
		Args: make([]sqlbuilder.Arg, 0, cap),
	}
}

func (w *ArgWriter) WriteArg(arg sqlbuilder.Arg) error {
	w.Args = append(w.Args, arg)
	return nil
}

func (b *fixedBuilder) WriteString(s string) (n int, err error) {
	if len(s)+b.size > b.cap {
		b.size = b.cap
		return b.builder.WriteString(s[0 : b.cap-b.size])
	}
	b.size += len(s)
	return b.builder.WriteString(s)
}

func (b *fixedBuilder) String() string {
	return b.builder.String()
}

func TestCustomClause(t *testing.T) {
	RegisterTestingT(t)
	t.Run("normal", func(t *testing.T) {
		c := sqlbuilder.NewCustomClause(func(sqlWriter io.StringWriter, argWriter sqlbuilder.ArgWriter, level int) error {
			sqlWriter.WriteString("abc")
			argWriter.WriteArg(1)
			argWriter.WriteArg("2")
			return nil
		})
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "abc"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("on error", func(t *testing.T) {
		c := sqlbuilder.NewCustomClause(func(sqlWriter io.StringWriter, argWriter sqlbuilder.ArgWriter, level int) error {
			return errors.New("demo error")
		})
		sqlWriter := newFixedBuilder(2)
		var argWriter = NewArgWriter(0)
		err := c.Parse(sqlWriter, argWriter, 0)
		Expect(err).ShouldNot(Succeed())
	})
}

func TestSimpleClause(t *testing.T) {
	RegisterTestingT(t)
	t.Run("auto new line", func(t *testing.T) {
		c := sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "abc", 1, "2")
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "abc\n"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("don't new line", func(t *testing.T) {
		c := sqlbuilder.NewSimpleClause(sqlbuilder.DontNewline, "abc", 1, "2")
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "abc"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("on error", func(t *testing.T) {
		c := sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "abc", 1, "2")
		sqlWriter := newFixedBuilder(2)
		var argWriter = NewArgWriter(0)
		err := c.Parse(sqlWriter, argWriter, 0)
		Expect(err).ShouldNot(Succeed())
	})
}

func TestAddLevel(t *testing.T) {
	RegisterTestingT(t)
	t.Run("auto new line", func(t *testing.T) {
		c := sqlbuilder.AddClauseLevel(sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "abc", 1, "2"), 1)
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "  abc\n"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("don't new line", func(t *testing.T) {
		c := sqlbuilder.AddClauseLevel(sqlbuilder.NewSimpleClause(sqlbuilder.DontNewline, "abc", 1, "2"), 1)
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "  abc"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("change to compact", func(t *testing.T) {
		c := sqlbuilder.AddClauseLevel(sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "abc", 1, "2"), sqlbuilder.Compact)
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "abc "
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("on error", func(t *testing.T) {
		c := sqlbuilder.AddClauseLevel(sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "abc", 1, "2"), 1)
		sqlWriter := newFixedBuilder(2)
		var argWriter = NewArgWriter(0)
		err := c.Parse(sqlWriter, argWriter, 0)
		Expect(err).ShouldNot(Succeed())
	})
}

func TestGetSpace(t *testing.T) {
	RegisterTestingT(t)
	const singleSpace = "  " // Same as the singleSpace in sqlbuilder/clause.go
	t.Run("compact", func(t *testing.T) {
		Expect(sqlbuilder.GetSpace(sqlbuilder.Compact)).Should(Equal(""))
	})
	t.Run("one space", func(t *testing.T) {
		Expect(sqlbuilder.GetSpace(1)).Should(Equal(singleSpace))
	})
	t.Run("two spaces", func(t *testing.T) {
		Expect(sqlbuilder.GetSpace(2)).Should(Equal(singleSpace + singleSpace))
	})
	t.Run("optimized spaces", func(t *testing.T) {
		Expect(sqlbuilder.GetSpace(31)).Should(Equal(strings.Repeat(singleSpace, 31)))
	})
	t.Run("unoptimized spaces", func(t *testing.T) {
		Expect(sqlbuilder.GetSpace(32)).Should(Equal(strings.Repeat(singleSpace, 32)))
	})
}

func TestWriteString(t *testing.T) {
	RegisterTestingT(t)
	t.Run("write susccess", func(t *testing.T) {
		sqlWriter := newFixedBuilder(3)
		err := sqlbuilder.WriteString(sqlWriter, "123")
		Expect(err).Should(Succeed())
		Expect(sqlWriter.String()).Should(Equal("123"))
	})
	t.Run("write failure", func(t *testing.T) {
		sqlWriter := newFixedBuilder(2)
		err := sqlbuilder.WriteString(sqlWriter, "123")
		Expect(err).ShouldNot(Succeed())
	})
	t.Run("write twice susccess", func(t *testing.T) {
		var err error
		sqlWriter := newFixedBuilder(6)
		err = sqlbuilder.WriteString(sqlWriter, "123")
		Expect(err).Should(Succeed())
		err = sqlbuilder.WriteString(sqlWriter, "456")
		Expect(err).Should(Succeed())
		Expect(sqlWriter.String()).Should(Equal("123456"))
	})
	t.Run("write twice failure", func(t *testing.T) {
		var err error
		sqlWriter := newFixedBuilder(3)
		err = sqlbuilder.WriteString(sqlWriter, "123")
		Expect(err).Should(Succeed())
		err = sqlbuilder.WriteString(sqlWriter, "456")
		Expect(err).ShouldNot(Succeed())
	})
}

func TestWriteStringWithSpace(t *testing.T) {
	RegisterTestingT(t)
	t.Run("write without space susccess", func(t *testing.T) {
		sqlWriter := newFixedBuilder(3)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", sqlbuilder.Format)
		Expect(err).Should(Succeed())
		Expect(sqlWriter.String()).Should(Equal("123"))
	})
	t.Run("write without space failure", func(t *testing.T) {
		sqlWriter := newFixedBuilder(2)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", sqlbuilder.Format)
		Expect(err).ShouldNot(Succeed())
	})

	t.Run("write compact susccess", func(t *testing.T) {
		sqlWriter := newFixedBuilder(3)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", sqlbuilder.Compact)
		Expect(err).Should(Succeed())
		Expect(sqlWriter.String()).Should(Equal("123"))
	})
	t.Run("write compact space failure", func(t *testing.T) {
		sqlWriter := newFixedBuilder(2)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", sqlbuilder.Compact)
		Expect(err).ShouldNot(Succeed())
	})
	t.Run("write with space susccess", func(t *testing.T) {
		sqlWriter := newFixedBuilder(5)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", 1)
		Expect(err).Should(Succeed())
		Expect(sqlWriter.String()).Should(Equal("  123"))
	})
	t.Run("write with space failure 1", func(t *testing.T) {
		sqlWriter := newFixedBuilder(1)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", 1)
		Expect(err).ShouldNot(Succeed())
	})
	t.Run("write with space failure 2", func(t *testing.T) {
		sqlWriter := newFixedBuilder(2)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", 1)
		Expect(err).ShouldNot(Succeed())
	})
	t.Run("write with many space susccess", func(t *testing.T) {
		sqlWriter := newFixedBuilder(7)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", 2)
		Expect(err).Should(Succeed())
		Expect(sqlWriter.String()).Should(Equal("    123"))
	})
	t.Run("write with many space failure", func(t *testing.T) {
		sqlWriter := newFixedBuilder(4)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", 2)
		Expect(err).ShouldNot(Succeed())
	})
	t.Run("write with more space susccess", func(t *testing.T) {
		sqlWriter := newFixedBuilder(67)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", 32)
		Expect(err).Should(Succeed())
		Expect(sqlWriter.String()).Should(Equal(strings.Repeat("  ", 32) + "123"))
	})
	t.Run("write with more space failure", func(t *testing.T) {
		sqlWriter := newFixedBuilder(1)
		err := sqlbuilder.WriteStringWithSpace(sqlWriter, "123", 32)
		Expect(err).ShouldNot(Succeed())
	})
}
