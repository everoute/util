package sqlbuilder_test

import (
	"bytes"
	"io"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/everoute/util/sql/sqlbuilder"
)

func TestSimpleCondition(t *testing.T) {
	RegisterTestingT(t)
	c := sqlbuilder.NewCondition("col1=@col1", 1)
	buff := bytes.NewBufferString("")
	var argWriter = NewArgWriter(0)
	err := c.Parse(buff, argWriter)
	Expect(err).Should(Succeed())
	res := buff.String()
	ept := "col1=@col1"
	Expect(res).To(Equal(ept))
	resArgs := argWriter.Args
	eptArgs := []sqlbuilder.Arg{1}
	Expect(resArgs).To(Equal(eptArgs))
}

func TestCustomCondition(t *testing.T) {
	RegisterTestingT(t)
	c := sqlbuilder.NewCustomCondition(func(sqlWriter io.StringWriter, argWriter sqlbuilder.ArgWriter) error {
		var err error
		err = sqlbuilder.WriteString(sqlWriter, "col='@col'")
		if err != nil {
			return err
		}
		err = sqlbuilder.WriteArgs(argWriter, "OwO")
		if err != nil {
			return err
		}
		return nil
	})
	buff := bytes.NewBufferString("")
	var argWriter = NewArgWriter(0)
	err := c.Parse(buff, argWriter)
	Expect(err).Should(Succeed())
	res := buff.String()
	ept := "col='@col'"
	Expect(res).To(Equal(ept))
	resArgs := argWriter.Args
	eptArgs := []sqlbuilder.Arg{"OwO"}
	Expect(resArgs).To(Equal(eptArgs))
}

func TestAndCondition(t *testing.T) {
	RegisterTestingT(t)
	t.Run("save brackets", func(t *testing.T) {
		c1 := sqlbuilder.NewCondition("col1=@col1", 1)
		c2 := sqlbuilder.NewCondition("col2=@col2", "2")
		c := sqlbuilder.And(c1, c2, sqlbuilder.SaveBrackets)
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "(col1=@col1 AND col2=@col2)"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("omit brackets", func(t *testing.T) {
		c1 := sqlbuilder.NewCondition("col1=@col1", 1)
		c2 := sqlbuilder.NewCondition("col2=@col2", "2")
		c := sqlbuilder.And(c1, c2, sqlbuilder.OmitBrackets)
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "col1=@col1 AND col2=@col2"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
}

func TestOrCondition(t *testing.T) {
	RegisterTestingT(t)
	t.Run("save brackets", func(t *testing.T) {
		c1 := sqlbuilder.NewCondition("col1=@col1", 1)
		c2 := sqlbuilder.NewCondition("col2=@col2", "2")
		c := sqlbuilder.Or(c1, c2, sqlbuilder.SaveBrackets)
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "(col1=@col1 OR col2=@col2)"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("omit brackets", func(t *testing.T) {
		c1 := sqlbuilder.NewCondition("col1=@col1", 1)
		c2 := sqlbuilder.NewCondition("col2=@col2", "2")
		c := sqlbuilder.Or(c1, c2, sqlbuilder.OmitBrackets)
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "col1=@col1 OR col2=@col2"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2"}
		Expect(resArgs).To(Equal(eptArgs))
	})
}

func TestNotCondition(t *testing.T) {
	RegisterTestingT(t)
	t.Run("save brackets", func(t *testing.T) {
		sub := sqlbuilder.NewCondition("col=@col", 1)
		c := sqlbuilder.Not(sub, sqlbuilder.SaveBrackets)
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "(NOT col=@col)"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("omit brackets", func(t *testing.T) {
		sub := sqlbuilder.NewCondition("col=@col", 1)
		c := sqlbuilder.Not(sub, sqlbuilder.OmitBrackets)
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "NOT col=@col"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1}
		Expect(resArgs).To(Equal(eptArgs))
	})
}
