package pool_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/everoute/util/pool"
)

func TestStdPool(t *testing.T) {
	RegisterTestingT(t)

	type data4096 struct {
		_ [4096]byte
	}

	start := make(chan struct{})

	customChan := make(chan *data4096, 1000)
	var allocCount uint32 = 0
	p := pool.NewStdPoll[data4096](func() any {
		atomic.AddUint32(&allocCount, 1)
		return new(data4096)
	})
	var putter pool.Putter[data4096] = p
	var getter pool.Getter[data4096] = p

	for i := 0; i < 100; i++ {
		go func() {
			<-start
			for i := 0; i < 10000; i++ {
				customChan <- getter.Get()
			}
		}()
	}

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			<-start
			for i := 0; i < 100000; i++ {
				ptr := <-customChan
				putter.Put(ptr)
			}
			wg.Done()
		}()
	}

	close(start)
	wg.Wait()
	fmt.Println("alloc:", allocCount)
}
