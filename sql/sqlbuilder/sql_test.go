package sqlbuilder_test

import (
	"github.com/everoute/util/sql/sqlbuilder"
)

type ArgWriter struct {
	Args []sqlbuilder.Arg
}

func (w *ArgWriter) WriteArg(arg sqlbuilder.Arg) error {
	w.Args = append(w.Args, arg)
	return nil
}
