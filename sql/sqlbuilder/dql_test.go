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
	t.Run("name AS table", func(t *testing.T) {
		var c = sqlbuilder.WithClause{
			Tables: []sqlbuilder.NamedTable{
				{
					"a",
					sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.A"),
				},
				{
					"b",
					sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.B"),
				},
			},
			NamePosition: sqlbuilder.NameFirst,
		}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "WITH\na AS (\n  SELECT * FROM demo.A\n),\nb AS (\n  SELECT * FROM demo.B\n)\n"
		klog.Infof("sql expect:\n%s\nresult:\n%s\n", ept, res)
		Expect(res).To(Equal(ept))
	})
	t.Run("table AS name", func(t *testing.T) {
		var c = sqlbuilder.WithClause{
			Tables: []sqlbuilder.NamedTable{
				{
					"a",
					sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.A"),
				},
				{
					"b",
					sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.B"),
				},
			},
			NamePosition: sqlbuilder.NameAfter,
		}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "WITH\n(\n  SELECT * FROM demo.A\n) AS a,\n(\n  SELECT * FROM demo.B\n) AS b\n"
		klog.Infof("sql expect:\n%s\nresult:\n%s\n", ept, res)
		Expect(res).To(Equal(ept))
	})
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
	var argWriter = NewArgWriter(0)
	err := c.Parse(buff, argWriter, 0)
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

func TestFrom(t *testing.T) {
	RegisterTestingT(t)
	t.Run("from table name", func(t *testing.T) {
		c := sqlbuilder.From{
			Name: "demo_table",
		}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "FROM demo_table\n"
		klog.Infof("expect:\n%s\nresult:\n%s\n", ept, res)
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := make([]sqlbuilder.Arg, 0)
		klog.Infof("arg expect:\n%+v\nresult:\n%+v\n", eptArgs, resArgs)
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("from sub query", func(t *testing.T) {
		dql := sqlbuilder.DQL{
			Select: sqlbuilder.Select{
				Columns: []string{
					"max(x) AS max_x",
					"max(y) AS max_y",
				},
			},
			From: sqlbuilder.From{
				Name: "demo_table",
			},
		}
		c := sqlbuilder.From{
			Table: &dql,
		}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "FROM (\n  SELECT\n    max(x) AS max_x,\n    max(y) AS max_y\n  FROM demo_table\n)\n"
		klog.Infof("expect:\n%s\nresult:\n%s\n", ept, res)
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := make([]sqlbuilder.Arg, 0)
		klog.Infof("arg expect:\n%+v\nresult:\n%+v\n", eptArgs, resArgs)
		Expect(resArgs).To(Equal(eptArgs))
	})
}

func TestDQL(t *testing.T) {
	RegisterTestingT(t)
	t.Run("with space", func(t *testing.T) {
		dql := sqlbuilder.DQL{
			With: &sqlbuilder.WithClause{
				Tables: []sqlbuilder.NamedTable{
					{
						"a",
						sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.A"),
					},
					{
						"b",
						sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.B"),
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
				sqlbuilder.NewSimpleClause(sqlbuilder.DontNewline, "x >= 2"),
				sqlbuilder.NewSimpleClause(sqlbuilder.DontNewline, "y != 'a'"),
			},
			Group: sqlbuilder.MakeGroupby("x"),
			Order: sqlbuilder.MakeOrderby("y"),
			Limit: sqlbuilder.MakeLimit("1"),
			Additional: []sqlbuilder.Clause{
				sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "QWERTY"),
				sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "asdfgh"),
			},
		}
		buff := bytes.NewBufferString("")
		err := dql.Parse(buff, nil, sqlbuilder.Format)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "WITH\na AS (\n  SELECT * FROM demo.A\n),\nb AS (\n  SELECT * FROM demo.B\n)\nSELECT\n  max(x) AS max_x,\n  max(y) AS max_y\nFROM demo\nWHERE\n  x >= 2\n  AND y != 'a'\nGROUP BY x\nORDER BY y\nLIMIT 1\nQWERTY\nasdfgh\n"
		klog.Infof("expect:\n%s\nresult:\n%s\n", ept, res)
		Expect(res).To(Equal(ept))
	})
	t.Run("without space", func(t *testing.T) {
		dql := sqlbuilder.DQL{
			With: &sqlbuilder.WithClause{
				Tables: []sqlbuilder.NamedTable{
					{
						"a",
						sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.A"),
					},
					{
						"b",
						sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.B"),
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
				sqlbuilder.NewSimpleClause(sqlbuilder.DontNewline, "x >= 2"),
				sqlbuilder.NewSimpleClause(sqlbuilder.DontNewline, "y != 'a'"),
			},
			Group: sqlbuilder.MakeGroupby("x"),
			Order: sqlbuilder.MakeOrderby("y"),
			Limit: sqlbuilder.MakeLimit("1"),
			Additional: []sqlbuilder.Clause{
				sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "QWERTY"),
				sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "asdfgh"),
			},
		}
		buff := bytes.NewBufferString("")
		err := dql.Parse(buff, nil, sqlbuilder.Compact)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "WITH a AS ( SELECT * FROM demo.A ), b AS ( SELECT * FROM demo.B ) SELECT max(x) AS max_x, max(y) AS max_y FROM demo WHERE x >= 2 AND y != 'a' GROUP BY x ORDER BY y LIMIT 1 QWERTY asdfgh "
		klog.Infof("expect:\n%s\nresult:\n%s\n", ept, res)
		Expect(res).To(Equal(ept))
	})
}
