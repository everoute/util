package sqlbuilder_test

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/everoute/util/sql/sqlbuilder"
)

func TestWith(t *testing.T) {
	t.Run("name AS table", func(t *testing.T) {
		RegisterTestingT(t)
		var c = sqlbuilder.WithClause{
			Tables: []sqlbuilder.Table{
				sqlbuilder.NameAsTable(
					"a",
					sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.A"),
				),
				sqlbuilder.NameAsTable(
					"b",
					sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.B"),
				),
			},
		}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "WITH\na AS (\n  SELECT * FROM demo.A\n),\nb AS (\n  SELECT * FROM demo.B\n)\n"
		Expect(res).To(Equal(ept))
	})
	t.Run("table AS name", func(t *testing.T) {
		RegisterTestingT(t)
		var c = sqlbuilder.WithClause{
			Tables: []sqlbuilder.Table{
				sqlbuilder.TableAsName(
					sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.A"),
					"a",
				),
				sqlbuilder.TableAsName(
					sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.B"),
					"b",
				),
			},
		}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "WITH\n(\n  SELECT * FROM demo.A\n) AS a,\n(\n  SELECT * FROM demo.B\n) AS b\n"
		Expect(res).To(Equal(ept))
	})
	t.Run("name only", func(t *testing.T) {
		RegisterTestingT(t)
		var c = sqlbuilder.WithClause{
			Tables: []sqlbuilder.Table{
				sqlbuilder.TableByName("a"),
				sqlbuilder.TableByName("b"),
			},
		}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "WITH\na,\nb\n"
		Expect(res).To(Equal(ept))
	})
}

func TestSelect(t *testing.T) {
	t.Run("SELECT *", func(t *testing.T) {
		RegisterTestingT(t)
		c := sqlbuilder.Select{}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "SELECT *\n"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("SELECT ...", func(t *testing.T) {
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
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{"x", 2}
		Expect(resArgs).To(Equal(eptArgs))
	})
}

func TestFrom(t *testing.T) {
	t.Run("from table name", func(t *testing.T) {
		RegisterTestingT(t)
		c := sqlbuilder.From{sqlbuilder.TableByName("demo_table")}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "FROM demo_table\n"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := make([]sqlbuilder.Arg, 0)
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("from sub query", func(t *testing.T) {
		RegisterTestingT(t)
		dql := sqlbuilder.DQL{
			Select: sqlbuilder.Select{
				Columns: []string{
					"max(x) AS max_x",
					"max(y) AS max_y",
				},
			},
			From: sqlbuilder.From{sqlbuilder.TableByName("demo_table")},
		}
		c := sqlbuilder.From{sqlbuilder.TableByClause(&dql)}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "FROM (\n  SELECT\n    max(x) AS max_x,\n    max(y) AS max_y\n  FROM demo_table\n)\n"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := make([]sqlbuilder.Arg, 0)
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("from name AS table", func(t *testing.T) {
		RegisterTestingT(t)
		dql := sqlbuilder.DQL{
			Select: sqlbuilder.Select{
				Columns: []string{
					"max(x) AS max_x",
					"max(y) AS max_y",
				},
			},
			From: sqlbuilder.From{sqlbuilder.TableByName("demo_table")},
		}
		c := sqlbuilder.From{sqlbuilder.NameAsTable("another_table", &dql)}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "FROM another_table AS (\n  SELECT\n    max(x) AS max_x,\n    max(y) AS max_y\n  FROM demo_table\n)\n"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := make([]sqlbuilder.Arg, 0)
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("from table AS name", func(t *testing.T) {
		RegisterTestingT(t)
		dql := sqlbuilder.DQL{
			Select: sqlbuilder.Select{
				Columns: []string{
					"max(x) AS max_x",
					"max(y) AS max_y",
				},
			},
			From: sqlbuilder.From{sqlbuilder.TableByName("demo_table")},
		}
		c := sqlbuilder.From{sqlbuilder.TableAsName(&dql, "another_table")}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := c.Parse(buff, argWriter, 0)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "FROM (\n  SELECT\n    max(x) AS max_x,\n    max(y) AS max_y\n  FROM demo_table\n) AS another_table\n"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := make([]sqlbuilder.Arg, 0)
		Expect(resArgs).To(Equal(eptArgs))
	})
}

func TestDQL(t *testing.T) {
	t.Run("with space", func(t *testing.T) {
		RegisterTestingT(t)
		dql := sqlbuilder.DQL{
			With: &sqlbuilder.WithClause{
				Tables: []sqlbuilder.Table{
					sqlbuilder.NameAsTable(
						"a",
						sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.A"),
					),
					sqlbuilder.NameAsTable(
						"b",
						sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.B"),
					),
				},
			},
			Select: sqlbuilder.Select{
				Columns: []string{
					"max(x) AS max_x",
					"max(y) AS max_y",
				},
			},
			From: sqlbuilder.From{sqlbuilder.TableByName("demo_table")},
			Where: sqlbuilder.WhereClause{
				[]sqlbuilder.Condition{
					sqlbuilder.NewCondition("x >= @x", 1),
					sqlbuilder.NewCondition("y != @y", "2"),
				},
			},
			Group: sqlbuilder.MakeGroupby("x"),
			Having: sqlbuilder.HavingClause{
				[]sqlbuilder.Condition{
					sqlbuilder.NewCondition("z < @z", 3.0),
				},
			},
			Order: sqlbuilder.MakeOrderby("y"),
			Limit: sqlbuilder.MakeLimit("1"),
			Additional: []sqlbuilder.Clause{
				sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "QWERTY"),
				sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "asdfgh"),
			},
		}
		buff := bytes.NewBufferString("")
		var argWriter = NewArgWriter(0)
		err := dql.Parse(buff, argWriter, sqlbuilder.Format)
		Expect(err).Should(Succeed())
		res := buff.String()
		ept := "WITH\na AS (\n  SELECT * FROM demo.A\n),\nb AS (\n  SELECT * FROM demo.B\n)\nSELECT\n  max(x) AS max_x,\n  max(y) AS max_y\nFROM demo_table\nWHERE\n  x >= @x\n  AND y != @y\nGROUP BY x\nHAVING\n  z < @z\nORDER BY y\nLIMIT 1\nQWERTY\nasdfgh\n"
		Expect(res).To(Equal(ept))
		resArgs := argWriter.Args
		eptArgs := []sqlbuilder.Arg{1, "2", 3.0}
		Expect(resArgs).To(Equal(eptArgs))
	})
	t.Run("without space", func(t *testing.T) {
		RegisterTestingT(t)
		dql := sqlbuilder.DQL{
			With: &sqlbuilder.WithClause{
				Tables: []sqlbuilder.Table{
					sqlbuilder.NameAsTable(
						"a",
						sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.A"),
					),
					sqlbuilder.NameAsTable(
						"b",
						sqlbuilder.NewSimpleClause(sqlbuilder.AutoNewline, "SELECT * FROM demo.B"),
					),
				},
			},
			Select: sqlbuilder.Select{
				Columns: []string{
					"max(x) AS max_x",
					"max(y) AS max_y",
				},
			},
			From: sqlbuilder.From{sqlbuilder.TableByName("demo_table")},
			Where: sqlbuilder.WhereClause{
				[]sqlbuilder.Condition{
					sqlbuilder.NewCondition("x >= 2"),
					sqlbuilder.NewCondition("y != 'a'"),
				},
			},
			Group: sqlbuilder.MakeGroupby("x"),
			Having: sqlbuilder.HavingClause{
				[]sqlbuilder.Condition{
					sqlbuilder.NewCondition("z < @z"),
				},
			},
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
		ept := "WITH a AS ( SELECT * FROM demo.A ), b AS ( SELECT * FROM demo.B ) SELECT max(x) AS max_x, max(y) AS max_y FROM demo_table WHERE x >= 2 AND y != 'a' GROUP BY x HAVING z < @z ORDER BY y LIMIT 1 QWERTY asdfgh "
		Expect(res).To(Equal(ept))
	})
}
