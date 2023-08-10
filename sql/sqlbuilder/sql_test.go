package sqlbuilder_test

import (
	"github.com/everoute/util/sql/sqlbuilder"
)

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
