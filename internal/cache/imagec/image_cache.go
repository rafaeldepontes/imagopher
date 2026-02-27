package imagec

import (
	"runtime"
	"time"

	"github.com/rafaeldepontes/imagopher/internal/cache"
)

type data[T any] struct {
	value     T
	expiresAt time.Time
}

type imgCache[T any] struct {
	memory map[string]*data[T]
}

func NewCache[T any]() cache.Cache[string, T] {
	return &imgCache[T]{
		memory: make(map[string]*data[T]),
	}
}

func (i *imgCache[T]) Add(key string, value T) {
	i.AddWithTLS(key, value, nil)
}

// Add implements [cache.Cache].
func (i *imgCache[T]) AddWithTLS(key string, value T, duration *time.Duration) {
	expiresAt := time.Now()

	var defaultDuration *time.Duration = new((time.Duration)(cache.DefaultDuration))
	if duration == nil {
		expiresAt = expiresAt.Add(*defaultDuration * time.Minute)
	} else {
		expiresAt = expiresAt.Add(*duration * time.Minute)
	}

	i.memory[key] = &data[T]{
		value:     value,
		expiresAt: expiresAt,
	}
}

// Clear implements [cache.Cache].
func (i *imgCache[T]) Clear() {
	clear(i.memory)
}

// FullClear implements [cache.Cache].
func (i *imgCache[T]) FullClear() {
	clear(i.memory)
	runtime.GC()
}

// Get implements [cache.Cache].
func (i *imgCache[T]) Get(key string) (T, bool) {
	data, has := i.memory[key]
	if !has {
		return zeroVal[T](), false
	}

	if data.expiresAt.Before(time.Now()) {
		i.Remove(key)
		return zeroVal[T](), false
	}

	return data.value, true
}

// Remove implements [cache.Cache].
func (i *imgCache[T]) Remove(key string) {
	delete(i.memory, key)
}

// Set implements [cache.Cache].
func (i *imgCache[T]) Set(key string, value T) {
	data, has := i.memory[key]
	if !has {
		i.Add(key, value)
		return
	}

	if data.expiresAt.Before(time.Now()) {
		i.Remove(key)
		i.Add(key, value)
		return
	}

	data.value = value
	data.expiresAt = time.Now().Add(cache.DefaultDuration)
}

func zeroVal[T any]() T {
	var zeroVal T
	return zeroVal
}
