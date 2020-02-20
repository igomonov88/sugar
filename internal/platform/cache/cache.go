package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// ErrInvalidConfig is used then some of the config values does not specified
var ErrInvalidConfig = errors.New("config values does not specified properly")

// Config store required properties to use Cache
type Config struct {
	// DefaultDuration represents the default value of TimeToLeave parameter of any item in cache. If TTL of the item is
	// expired we delete the item from the cache.
	DefaultDuration time.Duration
	// Size represents max size of the cache,excess of which entails evict of the item from the cache with LRU mechanics
	Size int
}

// Cache is the representation of the LRUCache item
type Cache struct {
	lock            sync.Mutex
	entryList       *list.List
	items           map[string]*list.Element
	defaultDuration time.Duration
	initialSize     int
	currentSize     int
}

// entry is a value which will store key/value pair in the cache
type entry struct {
	key       string
	value     interface{}
	expiresAt time.Time
}

// NewCache knows how to create Cache with given configuration
func NewCache(cfg Config) (*Cache, error) {
	if cfg.Size == 0 || cfg.DefaultDuration == 0 {
		return nil, ErrInvalidConfig
	}
	return &Cache{
		entryList:       list.New(),
		items:           make(map[string]*list.Element, cfg.Size),
		lock:            sync.Mutex{},
		defaultDuration: cfg.DefaultDuration,
		initialSize:     cfg.Size,
	}, nil
}

// Add knows how to add value or values to cache with given key
func (c *Cache) Add(key string, value interface{}) {
	c.lock.Lock()
	add(c, key, value)
	c.lock.Unlock()
}

// add adds element to list if in this process we going out of the initial cache size we will remove appropriate item
// from the cache
func add(cache *Cache, key string, value interface{}) {
	if cache.currentSize >= cache.initialSize {
		remove(cache)
	}
	if element, ok := cache.items[key]; ok && value != nil {
		cache.entryList.MoveToFront(element)
		entry := element.Value.(*entry)
		entry.value = value
		entry.expiresAt = time.Now().Add(cache.defaultDuration)
		return
	}
	element := cache.entryList.PushFront(&entry{key:key,value:value, expiresAt:time.Now().Add(cache.defaultDuration)})
	cache.items[key]= element
	cache.currentSize++
}

// Get knows hot to get value from the cache, and if it exists in cache it returns cached value and true as a second
// parameter to be able to operate with cache in map style. It makes it in concurrent safe way, with using a mutex lock
func (c *Cache) Get(key string) (value interface{}, exist bool) {
	c.lock.Lock()
	value, exist = get(c, key)
	c.lock.Unlock()
	return
}

// get returns value from the cache if it exists there and appropriate bool value of the exist parameter
func get (cache *Cache, key string) (interface{}, bool) {
	if element, exist := cache.items[key]; exist {
		entry := element.Value.(*entry)
		// if entry time to live is expired, we delete this entry form the cache
		if time.Now().After(entry.expiresAt) {
			cache.entryList.Remove(element)
			delete(cache.items, key)
			if cache.currentSize > 0 {
				cache.currentSize--
			}
			return nil, false
		}
		cache.entryList.MoveToFront(element)
		return entry.value, exist
	}
	return nil, false
}

// Purge knows hot to purge cache
func (c *Cache)Purge() {
	c.lock.Lock()
	purge(c)
	c.lock.Unlock()
}

// purge making purging cache
func purge (cache *Cache) {
	cache.entryList = list.New()
	cache.items = make(map[string]*list.Element, cache.initialSize)
	cache.currentSize = 0
}

// remove removes element from the cache
func remove(cache *Cache) {
	element := cache.entryList.Back()
	delete(cache.items, element.Value.(*entry).key)
	cache.entryList.Remove(element)
	cache.currentSize--
}
