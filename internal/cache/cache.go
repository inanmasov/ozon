package cache

import (
	"errors"
	"fmt"
	"log"

	"github.com/patrickmn/go-cache"
)

type myCache struct {
	cache *cache.Cache
}

var GlobalCacheManager *myCache

func InitCache() *myCache {
	fmt.Println(GlobalCacheManager)
	if GlobalCacheManager == nil {
		GlobalCacheManager = &myCache{
			cache: cache.New(cache.NoExpiration, cache.NoExpiration),
		}
	}
	return GlobalCacheManager
}

func (mc *myCache) SaveURL(originalURL, shortURL string) error {
	// Проверка наличия значения в кеше
	for _, val := range mc.cache.Items() {
		if val.Object == originalURL {
			log.Println("Such data is already in the cache")
			return errors.New("such data is already in the cache")
		}
	}

	// Проверка значения в кеше
	_, found := mc.cache.Get(shortURL)
	if !found {
		mc.cache.Set(shortURL, originalURL, -1)
		return nil
	} else {
		log.Println("Such data is already in the cache")
		return errors.New("such data is already in the cache")
	}
}

func (mc *myCache) GetOriginalURL(shortURL string) (string, error) {
	// Получение значения в кеше
	val, found := mc.cache.Get(shortURL)
	if found {
		str, ok := val.(string)
		if ok {
			return str, nil
		} else {
			log.Println("Link value error")
			return "", errors.New("link value error")
		}

	} else {
		log.Println("There is no such data in the cachee")
		return "", errors.New("there is no such data in the cache")
	}
}

func (mc *myCache) CheckShortUrl(shortURL string) error {
	// Проверка значения в кеше
	_, found := mc.cache.Get(shortURL)
	if found {
		return errors.New("short URL is already in the cache")
	} else {
		return nil
	}
}
