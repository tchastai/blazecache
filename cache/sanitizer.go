package cache

import "time"

type sanitizer struct {
	interval time.Duration
	stop     chan bool
}

func newSanitizer(interval time.Duration) *sanitizer {
	return &sanitizer{
		interval: interval,
		stop:     make(chan bool),
	}
}

func (s *sanitizer) run(w *worker) {
	ticker := time.NewTicker(s.interval)
	for {
		select {
		case <-ticker.C:
			w.DeleteExpired()
		case <-s.stop:
			ticker.Stop()
			return
		}
	}
}
