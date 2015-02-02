// Package lock provides a read-write mutex implementation which supports
// non-blocking locking
package lock

// TryRWMutex is similar to RWMutex; it can be locked by a single writer or any
// number of readers. It only supports a non-blocking, best-effort write lock.
type TryRWMutex struct {
	ch chan struct{} // write lock
	rs uint64        // number of active readers
}

// TryLock attempts to take a write lock, returning a bool indicating whether
// the attempt was successful. It does not block.
func (t *TryRWMutex) TryLock() bool {
	if !t.intTryLock() {
		return false
	}
	if t.rs != 0 {
		t.Unlock()
		return false
	}
	return true
}

// TryRLock attempts to take a read lock, returning a bool indicating whether
// the attempt was successful. It does not block.
func (t *TryRWMutex) TryRLock() bool {
	if !t.TryLock() {
		return false
	}
	t.rs++
	t.Unlock()
	return true
}

// RLock is like TryRLock but will block until it succeeds.
func (t *TryRWMutex) RLock() {
	t.intLock()
	t.rs++
	t.Unlock()
}

// Unlock unlocks the Lock for writing. It must only be called
// while the lock is held, or it will block.
func (t *TryRWMutex) Unlock() {
	<-t.ch
}

// intTryLockLock attempts to take the internal lock in a non-blocking way,
// returning a bool indicating whether the attempt was successful.
func (t *TryRWMutex) intTryLock() bool {
	select {
	case t.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// intLock is like intTryLock but will block until it succeeds.
func (t *TryRWMutex) intLock() {
	t.ch <- struct{}{}
}
