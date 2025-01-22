package cache

import (
	"time"
)

var (
	DefaultExpiration time.Duration = 0
	NoExpiration      time.Duration = -1
)

type Cache struct {
	*worker
}

func New(expiration time.Duration, interval time.Duration) *Cache {
	worker := newWorker(expiration)
	if interval > 0 {
		s := newSanitizer(interval)
		worker.sanitizer = s
		go s.run(worker)
	}
	return &Cache{
		worker: worker,
	}
}
