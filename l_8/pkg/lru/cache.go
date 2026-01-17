package lru

import "container/list"

// Якщо наш кеш вже повний (ми досягли нашого capacity)
// то має видалитись той елемент, який ми до якого ми доступались (читали) найдавніше

type LruCache interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type Cache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

type node struct {
	id    string
	value string
}

func NewLruCache(capacity int) LruCache {
	return &Cache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *Cache) Put(key, value string) {

	if el, ok := c.cache[key]; ok {
		el.Value.(*node).value = value
		c.list.MoveToFront(c.cache[key])
		return
	}

	if c.list.Len() == c.capacity {
		item := c.list.Back()
		id := item.Value.(*node).id
		delete(c.cache, id)
		c.list.Remove(item)
	}

	listItem := &node{
		id:    key,
		value: value,
	}

	resource := c.list.PushFront(listItem)
	c.cache[key] = resource
}

func (c *Cache) Get(key string) (string, bool) {
	if _, ok := c.cache[key]; ok {
		c.list.MoveToFront(c.cache[key])
		return c.cache[key].Value.(*node).value, true
	}
	return "", false
}
