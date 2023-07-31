package expected_test

import (
	"errors"
	"testing"

	"github.com/everoute/util/expected"
	. "github.com/onsi/gomega"
)

func TestEmpty(t *testing.T) {
	RegisterTestingT(t)
	// You shoundn't do anything with an Expected[T] object from Empty[T]() at anytime,
	// unless you called SetValue or SetValue on it.
	e := expected.Empty[int]()
	_ = e.IsOk()  // dangerous
	_ = e.IsBad() // dangerous

	e1 := expected.Empty[int]()
	e1.SetValue(1)
	Expect(e1.IsOk()).Should(BeTrue())   // safe
	Expect(e1.IsBad()).Should(BeFalse()) // safe

	e2 := expected.Empty[int]()
	e2.SetError(errors.New("demo error"))
	Expect(e2.IsOk()).Should(BeFalse()) // safe
	Expect(e2.IsBad()).Should(BeTrue()) // safe

}

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

func TestPack(t *testing.T) {
	RegisterTestingT(t)
	ev := expected.Pack(1, nil)
	Expect(ev.Error()).Should(Succeed())
	ee := expected.Pack(1, errors.New("demo error"))
	Expect(ee.Error()).ShouldNot(Succeed())
}

func TestPackNonNilP(t *testing.T) {
	RegisterTestingT(t)
	i := 1
	p := &i
	var pNil *int = nil
	ev := expected.PackNonNilP(p, nil)
	Expect(ev.Error()).Should(Succeed())
	eNil := expected.PackNonNilP(pNil, nil)
	Expect(eNil.Error()).ShouldNot(Succeed())
	ee := expected.PackNonNilP(p, errors.New("demo error"))
	Expect(ee.Error()).ShouldNot(Succeed())
}

func TestNonNilP(t *testing.T) {
	RegisterTestingT(t)
	i := 1
	p := &i
	var pNil *int = nil
	ev := expected.NonNilP(p)
	Expect(ev.Error()).Should(Succeed())
	eNil := expected.NonNilP(pNil)
	Expect(eNil.Error()).ShouldNot(Succeed())
}

func TestWrap(t *testing.T) {
	RegisterTestingT(t)
	half := func(i int) (int, error) {
		if i%2 == 0 {
			return i / 2, nil
		}
		return -1, errors.New("i is not an even number")
	}
	wrapedHalf := expected.Wrap(half)
	Expect(wrapedHalf(2).Error()).Should(Succeed())
	r, ok := wrapedHalf(2).Value()
	Expect(ok).Should(BeTrue())
	Expect(r).Should(Equal(1))
	Expect(wrapedHalf(3).Error()).ShouldNot(Succeed())
}

func TestWrapNonNilP(t *testing.T) {
	RegisterTestingT(t)
	half := func(i *int) (*int, error) {
		if *i%2 == 0 {
			r := *i / 2
			return &r, nil
		}
		return nil, errors.New("i is not an even number")
	}
	wrapedHalf := expected.WrapNonNilP(half)
	var two = 2
	var three = 3
	Expect(wrapedHalf(&two).Error()).Should(Succeed())
	r, ok := wrapedHalf(&two).Value()
	Expect(ok).Should(BeTrue())
	Expect(*r).Should(Equal(1))
	Expect(wrapedHalf(&three).Error()).ShouldNot(Succeed())
}

func TestUnwrap(t *testing.T) {
	RegisterTestingT(t)
	wrapedHalf := func(i int) expected.Expected[int] {
		if i%2 == 0 {
			return expected.Value(i / 2)
		}
		return expected.Error[int](errors.New("i is not an even number"))
	}
	unwrapedHalf := expected.Unwrap(wrapedHalf)
	var err error
	_, err = unwrapedHalf(1)
	Expect(err).ShouldNot(Succeed())
	v, err := unwrapedHalf(2)
	Expect(v).Should(Equal(1))
	Expect(err).Should(Succeed())
}

func TestSetValue(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Error[int](errors.New("demo error"))
	Expect(e.Error()).ShouldNot(Succeed())
	e.SetValue(1)
	v, err := e.Get()
	Expect(v).Should(Equal(1))
	Expect(err).Should(Succeed())
}

func TestSetError(t *testing.T) {
	RegisterTestingT(t)
	e := expected.Value(1)
	Expect(e.Error()).Should(Succeed())
	e.SetError(errors.New("demo error"))
	Expect(e.Error()).ShouldNot(Succeed())
}

func TestGet(t *testing.T) {
	RegisterTestingT(t)
	var v int
	var err error
	v, err = expected.Value(1).Get()
	Expect(v).Should(Equal(1))
	Expect(err).Should(Succeed())
	v, err = expected.Error[int](errors.New("demo error")).Get()
	_ = v // dangerous
	Expect(err).ShouldNot(Succeed())
	Expect(err.Error()).Should(Equal("demo error"))
}

func TestIsOk(t *testing.T) {
	RegisterTestingT(t)
	ev := expected.Value(1)
	Expect(ev.IsOk()).Should(BeTrue())
	ee := expected.Error[int](errors.New("demo error"))
	Expect(ee.IsOk()).Should(BeFalse())
}

func TestIsBad(t *testing.T) {
	RegisterTestingT(t)
	ev := expected.Value(1)
	Expect(ev.IsBad()).Should(BeFalse())
	ee := expected.Error[int](errors.New("demo error"))
	Expect(ee.IsBad()).Should(BeTrue())
}

func TestAndThen(t *testing.T) {
	RegisterTestingT(t)
	e1 := expected.Value(1)
	e2 := e1.AndThen(func(i int) expected.Expected[int] {
		return expected.Value(i + 1)
	})
	var err error
	v, err := e2.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(2))
	e3 := expected.Error[int](errors.New("demo error"))
	e4 := e3.AndThen(func(i int) expected.Expected[int] {
		t.Fail()
		return expected.Empty[int]()
	})
	_, err = e4.Get()
	Expect(err).ShouldNot(Succeed())
}

func TestOrElse(t *testing.T) {
	RegisterTestingT(t)
	e1 := expected.Value(1)
	e2 := e1.OrElse(func(error) expected.Expected[int] {
		t.Fail()
		return expected.Value(2)
	})
	var err error
	v, err := e2.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(1))
	e3 := expected.Error[int](errors.New("demo error"))
	e4 := e3.OrElse(func(error) expected.Expected[int] {
		return expected.Value(2)
	})
	v, err = e4.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(2))
}

func TestOr(t *testing.T) {
	RegisterTestingT(t)
	e1 := expected.Value(1)
	e2 := e1.Or(expected.Value(2))
	var err error
	v, err := e2.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(1))
	e3 := expected.Error[int](errors.New("demo error"))
	e4 := e3.Or(expected.Value(2))
	v, err = e4.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(2))
}

func TestTransform(t *testing.T) {
	RegisterTestingT(t)
	e1 := expected.Value(1)
	e2 := e1.Transform(func(i int) int {
		return i + 1
	})
	var err error
	v, err := e2.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(2))
	e3 := expected.Error[int](errors.New("demo error"))
	e4 := e3.Transform(func(i int) int {
		t.Fail()
		return 0
	})
	_, err = e4.Get()
	Expect(err).ShouldNot(Succeed())

}

func TestTransformError(t *testing.T) {
	RegisterTestingT(t)

	e1 := expected.Value(1)
	e2 := e1.TransformError(func(error) error {
		t.Fail()
		return errors.New("another error")
	})
	var err error
	v, err := e2.Get()
	Expect(err).Should(Succeed())
	Expect(v).Should(Equal(1))
	e3 := expected.Error[int](errors.New("demo error"))
	e4 := e3.TransformError(func(error) error {
		return errors.New("another error")
	})
	_, err = e4.Get()
	Expect(err).ShouldNot(Succeed())
	Expect(err.Error()).Should(Equal("another error"))
}
