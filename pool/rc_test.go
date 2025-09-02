package pool_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/everoute/util/pool"
)

func TestRCPool(t *testing.T) {
	var alloced int
	var freed int
	custromPool := pool.NewCustomPool(
		func() int {
			alloced++
			return 0
		},
		func(_ int) {
			freed++
		},
	)
	RegisterTestingT(t)
	p := pool.NewRCPool(custromPool)
	rc := p.New()
	newRc := rc.Ref()
	newRc.Unref()
	rc.Unref()
	Expect(alloced).To(Equal(freed))
}
