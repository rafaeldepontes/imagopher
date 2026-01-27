package cache

import "time"

// NOTE: This is just in memory, I should be using Redis for this operation,
// but I'm to lazy to prepare everything and supposedly we won't have that many
// users. But if comes to that, I will make this some type of interface to easy
// switch everything... Or maybe not...

const DefaultDuration = 15

// Cache implements a internal cache system, every method is not case sensitive
// to avoid case miss... So "A" and "a" gives the same value.
type Cache[K comparable, T any] interface {

	// Add adds something to cache, with no TLS were specified it will use the
	// default value of 15 minutes.
	Add(key K, value T, time *time.Duration)

	// Set updates the cache value and also refresh the TLS, if the TLS were
	// expired then it removes the old value and create a new one...
	Set(key K, value T)

	// Remove removes...
	Remove(key K)

	// Get gets the value if any and also returns a boolean to check if it is
	// expired.
	Get(key K) (T, bool)

	// Clear clears cache...
	Clear()

	// FullClear should only be called in case of lack of memory, ideally this
	// would never be used since it calls the GC directly and can, depending on
	// the workload, slow down A LOT... Be careful.
	FullClear()
}
