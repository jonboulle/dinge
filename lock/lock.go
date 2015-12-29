// Copyright 2015 Jonathan Boulle
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
