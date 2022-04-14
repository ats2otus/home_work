package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.RWMutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	item, exists := lru.items[key]
	if !exists {
		lru.items[key] = lru.queue.PushFront(value)
		if lru.queue.Len() > lru.capacity {
			lru.queue.Remove(lru.queue.Back())
		}
	} else {
		item.Value = value
		lru.queue.MoveToFront(item)
	}
	return exists
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	item, ok := lru.items[key]
	if !ok {
		return nil, false
	}
	lru.queue.MoveToFront(item)
	return item.Value, true
}

func (lru *lruCache) Clear() {
	for k := range lru.items {
		delete(lru.items, k)
	}
	lru.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
