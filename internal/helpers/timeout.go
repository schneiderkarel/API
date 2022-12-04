package helpers

import "time"

const DefaultTimeout = 10 * time.Second

func WithTimeout(fn func(), timeout time.Duration) bool {
	t := time.NewTimer(timeout)
	fnFinishedCh := make(chan struct{}, 1)
	go func() {
		fn()
		fnFinishedCh <- struct{}{}
	}()
	select {
	case <-t.C:
		return false
	case <-fnFinishedCh:
		return true
	}
}
