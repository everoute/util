package sqlbuilder_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/everoute/util/sql/sqlbuilder"
	. "github.com/onsi/gomega"
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

func (b *fixedBuilder) WriteString(s string) (n int, err error) {
	if len(s)+b.size > b.cap {
		return b.builder.WriteString(s[0 : b.cap-b.size])
	}
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
	t.Run("on error", func(t *testing.T) {
		c := sqlbuilder.AddClauseLevel(sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "abc", 1, "2"), 1)
		sqlWriter := newFixedBuilder(2)
		var argWriter = NewArgWriter(0)
		err := c.Parse(sqlWriter, argWriter, 0)
		Expect(err).ShouldNot(Succeed())
	})
}
