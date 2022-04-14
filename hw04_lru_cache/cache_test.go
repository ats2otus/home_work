package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		// Write me
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func Test_Eviction(t *testing.T) {

	t.Run("simple", func(t *testing.T) {
		const capacity = 10

		c := lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}

		for i := 0; i < 100; i++ {
			c.Set(Key(fmt.Sprintf("key-%d", i)), i)
		}
		require.Equal(t, capacity, c.queue.Len())
	})

	t.Run("complex", func(t *testing.T) {
		const capacity = 3

		c := lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}

		c.Set("1", 1)
		c.Set("2", 2)
		c.Set("3", 3) // queue: 3, 2, 1

		c.Get("1") // queue: 1, 3, 2
		c.Get("2") // queue: 2, 1, 3

		c.Set("4", 4) // queue: 4, 2, 1

		require.Equal(t, 4, c.queue.Front().Value)
		require.Equal(t, 1, c.queue.Back().Value)
	})
}

func Test_Cache_Clear(t *testing.T) {
	const capacity = 10

	c := lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}

	for i := 0; i < 100; i++ {
		c.Set(Key(fmt.Sprintf("key-%d", i)), i)
	}
	c.Clear()
	require.Equal(t, 0, c.queue.Len(), "queue len")
	require.Equal(t, 0, len(c.items), "items len")
}
