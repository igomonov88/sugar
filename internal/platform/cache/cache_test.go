package cache

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestCache(t *testing.T) {
	t.Log("Given we starting to test cache:")
	{
		cfg := Config{
			DefaultDuration: 1 * time.Millisecond,
			Size:            2,
		}
		rand.Seed(time.Now().UnixNano())
		cache, err := New(cfg)
		if err != nil {
			t.Fatalf("\t%s\t Should be able to get new cache instance: %s .", failed, err)
		}
		t.Logf("\t%s\t Should be able to get new cache instance.", success)

		{
			putItem := rand.Int31n(10)
			cache.Add("key", putItem)
			getItem, exist := cache.Get("key")
			if !exist {
				t.Fatalf("\t%s\t Should be able to get just pushed item from the cache.", failed)
			}
			if getItem != putItem {
				t.Fatalf("\t%s\t Should be able to get the same value as was pushed.", failed)
			}
			t.Logf("\t%s\t Should be able to get the same value as was pushed.", success)
		}

		{
			cache.Purge()
			if _, exist := cache.Get("key"); exist {
				t.Fatalf("\t%s\t Should be able to purge cache.", failed)
			}
			t.Logf("\t%s\t Should be able to purge cache.", success)
		}

		{
			for i := 0; i <= cfg.Size; i++ {
				cache.Add(fmt.Sprintf("key%v", strconv.Itoa(i)), rand.Int31n(10))
			}
			if _, exist := cache.Get("key0"); exist {
				t.Fatalf("\t%s\t Should remove old value if we try to add new item out of cache size.", failed)
			}
			t.Logf("\t%s\t Should remove old value if we try to add new item out of cache size.", success)
			cache.Purge()
		}

		{
			for i := 0; i < cfg.Size; i++ {
				cache.Add(fmt.Sprintf("key%v",strconv.Itoa(i)), i)
			}
			if _, exist := cache.Get("key0"); !exist {
				t.Fatalf("\t%s\t Should be able to get value from cache.", failed)
			}
			cache.Add("key2", rand.Int31n(10))

			if _, exist := cache.Get("key1"); exist{
				t.Fatalf("\t%s\t Should remove value with low usage from cache.", failed)
			}
			_, zeroItemExist := cache.Get("key0")
			_, secondItemExist := cache.Get("key2")
			if zeroItemExist == false && secondItemExist == false {
				t.Fatalf("\t%s\t Should be able to get last recently used items from cache.", failed)
			}
			t.Logf("\t%s\t Should be able to get last recently used items from cache.", success)
			cache.Purge()
		}

		{
			keys := make([]string, 0, cfg.Size)
			for i := 0; i < cfg.Size; i++ {
				cache.Add(fmt.Sprintf("key%v",strconv.Itoa(i)), i)
				keys = append(keys, fmt.Sprintf("key%v",strconv.Itoa(i)))
			}
			time.Sleep(10 * time.Millisecond)
			for i := range keys {
				if _, exist := cache.Get(keys[i]); exist {
					t.Fatalf("\t%s\t Should not beign able to get item with expired TTL.", failed)
				}
			}
			t.Logf("\t%s\t Should not being able to get item with expired TTL.", success)
		}
	}
}
