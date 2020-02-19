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
	// ExpirationTime represents the default value of TimeToLeave parameter of any item in cache. If TTL of the item is
	// expired we delete the item from the cache.
	ExpirationTime int64
	// Size represents max size of the cache,excess of which entails evict of the item from the cache with LRU mechanics
	Size int
}

// Cache is the representation of the LRUCache item
type Cache struct {
	entryList         *list.List
	items             map[string]*list.Element
	lock              sync.RWMutex
	defaultExpiration int64
	initialSize       int
	currentSize       int
}

// entry is a value which will store key/value pair in the cache
type entry struct {
	key        string
	value      interface{}
	expiresAt time.Time
}

// NewCache knows how to create Cache with given configuration
func NewCache(cfg Config) (*Cache, error) {
	if cfg.Size == 0 || cfg.ExpirationTime == 0 {return nil, ErrInvalidConfig}
	return &Cache{
		entryList: list.New(),
		items: make(map[string]*list.Element, cfg.Size),
		lock: sync.RWMutex{},
		defaultExpiration: cfg.ExpirationTime,
		initialSize: cfg.Size,
	}, nil
}

// Add knows how to add value or values to cache with given key
func (c *Cache) Add(key string, value interface{}) {
	c.lock.Lock()

	c.lock.Unlock()
}

func add(cache *Cache, key string, value interface{}) {

}