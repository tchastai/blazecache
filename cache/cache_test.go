package cache

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New(DefaultExpiration, 0)

	c.Set("a", "toto", DefaultExpiration)
	r, found := c.Get("a")
	if !found && r == "toto" {
		t.Error("Test failed : item not found : ", r)
	}

	c.Set("b", "toto", DefaultExpiration)
	r, found = c.Get("b")
	if !found && r == "toto" {
		t.Error("Test failed : item not found : ", r)
	}

	c.Set("c", "toto", DefaultExpiration)
	r, e, found := c.GetWithExpiration("c")
	if !found && e.IsZero() && r != nil {
		t.Error("Test failed : item not found with expiration : ", r)
	}

	c.Delete("c")
	_, found = c.Get("c")
	if found {
		t.Error("Test failed : item have not been deleted")
	}

	c.DeleteAll()
	length := c.GetLength()
	if length != 0 {
		t.Error("Test failed : items have not been deleted")
	}

	c.Set("d", "toto", DefaultExpiration)
	err := c.Add("d", "toto2", DefaultExpiration)
	if err == nil {
		t.Error("Test failed : an item with the same key already exists")
	}

	err = c.Replace("d", "toto2", DefaultExpiration)
	if err != nil {
		t.Error("Test failed : an item exists with the same key")
	}
}

func TestCacheWithSanitizer(t *testing.T) {
	c := New(2*time.Second, 5*time.Second)

	c.Set("a", "toto", DefaultExpiration)
	c.Set("b", "toto", 3*time.Second)
	c.Set("c", "toto", 4*time.Second)
	c.Set("d", "toto", 8*time.Second)
	c.Set("e", "toto", NoExpiration)

	time.Sleep(6 * time.Second)
	_, found := c.Get("a")
	if found {
		t.Error("Test failed : item should have been deleted")
	}
	_, found = c.Get("b")
	if found {
		t.Error("Test failed : item should have been deleted")
	}
	_, found = c.Get("c")
	if found {
		t.Error("Test failed : item should have been deleted")
	}
	_, found = c.Get("d")
	if !found {
		t.Error("Test failed : item should not be deleted")
	}
	_, found = c.Get("e")
	if !found {
		t.Error("Test failed : item should not be deleted")
	}
}

func benchmarkCacheGetConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	c := New(exp, 0)
	c.Set("a", "toto", DefaultExpiration)
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < each; j++ {
				c.Get("a")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkCacheGetConcurrentNotExpiring(b *testing.B) {
	benchmarkCacheGetConcurrent(b, NoExpiration)
}

func BenchmarkCacheGetConcurrentExpiring(b *testing.B) {
	benchmarkCacheGetConcurrent(b, 5*time.Minute)
}
