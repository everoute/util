package sqlbuilder_test

import (
	"bytes"
	"testing"

	"github.com/everoute/util/sql/sqlbuilder"
	. "github.com/onsi/gomega"
	"k8s.io/klog"
)

func TestWith(t *testing.T) {
	RegisterTestingT(t)
	var c = sqlbuilder.WithClause{
		Tables: []sqlbuilder.NamedTable{
			{
				"a",
				sqlbuilder.NewSimpleClause("SELECT * FROM demo.A\n"),
			},
			{
				"b",
				sqlbuilder.NewSimpleClause("SELECT * FROM demo.B\n"),
			},
		},
		NamePosition: sqlbuilder.NameFirst,
	}
	buff := bytes.NewBufferString("")
	var argWriter ArgWriter
	err := c.Parse(buff, &argWriter, 0)
	Expect(err).Should(Succeed())
	res := buff.String()
	ept := "WITH\na AS (\n  SELECT * FROM demo.A\n),\nb AS (\n  SELECT * FROM demo.B\n)\n"
	klog.Infof("sql expect:\n%s\nresult:\n%s\n", ept, res)
	Expect(res).To(Equal(ept))
}

func TestSelect(t *testing.T) {
	RegisterTestingT(t)
	c := sqlbuilder.Select{
		Columns: []string{
			"?",
			"max(y+?) AS max_y",
		},
		Args: []sqlbuilder.Arg{
			"x",
			2,
		},
	}
	buff := bytes.NewBufferString("")
	var argWriter ArgWriter
	err := c.Parse(buff, &argWriter, 0)
	Expect(err).Should(Succeed())
	res := buff.String()
	ept := "SELECT\n  ?,\n  max(y+?) AS max_y\n"
	klog.Infof("expect:\n%s\nresult:\n%s\n", ept, res)
	Expect(res).To(Equal(ept))
	resArgs := argWriter.Args
	eptArgs := []sqlbuilder.Arg{"x", 2}
	klog.Infof("arg expect:\n%+v\nresult:\n%+v\n", eptArgs, resArgs)
	Expect(resArgs).To(Equal(eptArgs))
}

func TestDQL(t *testing.T) {
	RegisterTestingT(t)
	dql := sqlbuilder.DQL{
		With: &sqlbuilder.WithClause{
			Tables: []sqlbuilder.NamedTable{
				{
					"a",
					sqlbuilder.NewSimpleClause("SELECT * FROM demo.A\n"),
				},
				{
					"b",
					sqlbuilder.NewSimpleClause("SELECT * FROM demo.B\n"),
				},
			},
			NamePosition: sqlbuilder.NameFirst,
		},
		Select: sqlbuilder.Select{
			Columns: []string{
				"max(x) AS max_x",
				"max(y) AS max_y",
			},
		},
		From: sqlbuilder.From{
			Name: "demo",
		},
		Where: []sqlbuilder.Condition{
			sqlbuilder.NewSimpleClause("x >= 2"),
			sqlbuilder.NewSimpleClause("y != 'a'"),
		},
		Group: sqlbuilder.MakeGroupby(false, "x"),
		Order: sqlbuilder.MakeOrderby(false, "y"),
		Limit: sqlbuilder.MakeLimit(false, "1"),
	}
	buff := bytes.NewBufferString("")
	err := dql.Parse(buff, nil, 0)
	Expect(err).Should(Succeed())
	res := buff.String()
	ept := "WITH\na AS (\n  SELECT * FROM demo.A\n),\nb AS (\n  SELECT * FROM demo.B\n)\nSELECT\n  max(x) AS max_x,\n  max(y) AS max_y\nFROM demo\nWHERE\nx >= 2\n  AND y != 'a'\nGROUP BY x\nORDER BY y\nLIMIT 1\n"
	klog.Infof("expect:\n%s\nresult:\n%s\n", ept, res)
	Expect(res).To(Equal(ept))
}
