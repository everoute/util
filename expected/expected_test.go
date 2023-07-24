package expected_test

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/everoute/util/expected"
)

func TestValue(t *testing.T) {
	RegisterTestingT(t)
	var value int = 1
	e := expected.Value(value)
	v, err := e.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(value))
	v, ok := e.Value()
	Expect(ok).Should(BeTrue())
	Expect(v).Should(Equal(value))
	err = e.Error()
	Expect(err).Should(Succeed())
}

func TestError(t *testing.T) {
	RegisterTestingT(t)
	var demoErr = errors.New("demo error")
	e := expected.Error[int](demoErr)
	v, err := e.Get()
	Expect(v).Should(Equal(0))
	Expect(err).ShouldNot(Succeed())
	v, ok := e.Value()
	Expect(v).Should(Equal(0))
	Expect(ok).Should(BeFalse())
	err = e.Error()
	Expect(err).Should(Equal(demoErr))
}

func TestValueOr(t *testing.T) {
	RegisterTestingT(t)
	var v int
	v = expected.Value(1).ValueOr(2)
	Expect(v).Should(Equal(1))

	v = expected.Error[int](errors.New("demo error")).ValueOr(2)
	Expect(v).Should(Equal(2))
}

func TestAndThenWithValue(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Value(1)
	e = e.AndThen(func(i int) expected.Expected[int] {
		return expected.Value(i + 1)
	})
	v, err := e.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(1 + 1))
	e = e.AndThen(func(i int) expected.Expected[int] {
		return expected.Error[int](errors.New("demo error"))
	})
	v, err = e.Get()
	Expect(err).ShouldNot(Succeed())
	Expect(v).ShouldNot(Equal(1 + 1))
}

func TestAndThenWithError(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Error[int](errors.New("demo error"))
	e.AndThen(func(i int) expected.Expected[int] {
		return expected.Value(i + 1)
	})
	err := e.Error()
	Expect(err).ShouldNot(Succeed())
}

func TestOrElseWithValue(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Value(1)
	err := e.OrElse(func(err error) expected.Expected[int] {
		t.Fail()
		return expected.Error[int](err)
	}).Error()
	Expect(err).Should(Succeed())
}

func TestOrElseWithError(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Error[int](errors.New("demo error"))
	e = e.OrElse(func(err error) expected.Expected[int] {
		return expected.Error[int](errors.New("another error"))
	})
	err := e.Error()
	Expect(err).ShouldNot(Succeed())
	e = e.OrElse(func(err error) expected.Expected[int] {
		return expected.Value(1)
	})
	v, err := e.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(1))
}

func TestTransformWithValue(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Value(1)
	e = e.Transform(func(i int) int {
		return i + 1
	})
	v, err := e.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(1 + 1))
}

func TestTransformWithError(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Error[int](errors.New("demo error"))
	e.Transform(func(i int) int {
		return i + 1
	})
	err := e.Error()
	Expect(err).ShouldNot(Succeed())
}

func TestTransformErrorWithValue(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Value(1)
	err := e.TransformError(func(err error) error {
		t.Fail()
		return err
	}).Error()
	Expect(err).Should(Succeed())
}

func TestTransformErrorWithError(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Error[int](errors.New("demo error"))
	e = e.TransformError(func(err error) error {
		return errors.New("another error")
	})
	err := e.Error()
	Expect(err).ShouldNot(Succeed())
}
