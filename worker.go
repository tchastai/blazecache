package cache

import (
	"fmt"
	"sync"
	"time"
)

type item struct {
	expiration int64
	object     interface{}
}

func newItem(expiration int64, object interface{}) *item {
	return &item{
		expiration: expiration,
		object:     object,
	}
}

type worker struct {
	mu                sync.RWMutex
	defaultExpiration time.Duration
	items             map[string]item
	sanitizer         *sanitizer
}

func newWorker(defaultExpiration time.Duration) *worker {
	items := make(map[string]item)
	return &worker{
		defaultExpiration: defaultExpiration,
		items:             items,
	}
}

func (w *worker) Set(key string, objewt interface{}, expiration time.Duration) {
	var exp int64

	if expiration == DefaultExpiration {
		expiration = w.defaultExpiration
	}

	if expiration > DefaultExpiration {
		exp = time.Now().Add(expiration).UnixNano()
	}

	w.mu.Lock()

	w.items[key] = *newItem(exp, objewt)

	w.mu.Unlock()
}

func (w *worker) set(key string, object interface{}, expiration time.Duration) {
	var exp int64

	if expiration == DefaultExpiration {
		expiration = w.defaultExpiration
	}

	if expiration > DefaultExpiration {
		exp = time.Now().Add(expiration).UnixNano()
	}

	w.items[key] = *newItem(exp, object)

}

func (w *worker) Add(key string, object interface{}, expiration time.Duration) error {
	w.mu.Lock()

	_, found := w.get(key)
	if found {
		w.mu.Unlock()
		return fmt.Errorf("item already exists with key = %s", key)
	}

	w.set(key, object, expiration)

	w.mu.Unlock()

	return nil
}

func (w *worker) Replace(key string, object interface{}, expiration time.Duration) error {

	w.mu.Lock()

	_, found := w.get(key)
	if !found {
		w.mu.Unlock()
		return fmt.Errorf("item does not exists with key = %s", key)
	}

	w.set(key, object, expiration)

	w.mu.Unlock()

	return nil
}

func (w *worker) Get(key string) (interface{}, bool) {
	w.mu.Lock()

	item, found := w.items[key]

	w.mu.Unlock()

	if (item.expiration > 0 && time.Now().UnixNano() > item.expiration) || !found {
		return nil, false
	}

	return item.object, true
}

func (w *worker) get(key string) (interface{}, bool) {

	item, found := w.items[key]
	if !found {
		return fmt.Errorf("item not found"), false
	}

	if (item.expiration > 0 && time.Now().UnixNano() > item.expiration) || !found {
		return nil, false
	}

	return item.object, true
}

func (w *worker) GetWithExpiration(key string) (interface{}, time.Time, bool) {
	w.mu.Lock()

	item, found := w.items[key]

	w.mu.Unlock()

	if (item.expiration > 0 && time.Now().UnixNano() > item.expiration) || !found {
		return nil, time.Time{}, false
	}

	return item.object, time.Unix(0, item.expiration), true
}

func (w *worker) Delete(key string) {
	w.mu.Lock()

	delete(w.items, key)

	w.mu.Unlock()
}

func (w *worker) DeleteAll() {
	w.mu.Lock()

	w.items = map[string]item{}

	w.mu.Unlock()
}

func (w *worker) GetLength() int {
	w.mu.Lock()

	l := len(w.items)

	w.mu.Unlock()

	return l
}

func (w *worker) DeleteExpired() {
	w.mu.Lock()
	for k, v := range w.items {
		if v.expiration != 0 && v.expiration < time.Now().Unix() {
			delete(w.items, k)
		}
	}
	w.mu.Unlock()
}
