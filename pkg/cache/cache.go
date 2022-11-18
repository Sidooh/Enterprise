package cache

import (
	"encoding/json"
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"reflect"
	"time"
)

var GlobalCache Cache[string, interface{}]

type Cache[K comparable, V any] interface {
	Get(key K) *V
	Set(key K, value V, time time.Duration) *V
	GetString(key K) string
	Unmarshal(key K, to interface{}) error
}

type cache[K comparable, V any] struct {
	cache *ttlcache.Cache[K, V]
}

func (c *cache[K, V]) GetString(key K) string {
	value := c.Get(key)
	if value != nil {
		if reflect.TypeOf(*value).Name() == "string" {
			return fmt.Sprint(*value)
		}
		return interfaceToString(value)
	}

	return ""
}

func (c *cache[K, V]) Get(key K) *V {
	value := c.cache.Get(key)
	if value != nil && !value.IsExpired() {
		v := value.Value()
		return &v
	}

	return nil
}

func (c *cache[K, V]) Set(key K, value V, time time.Duration) *V {
	v := c.cache.Set(key, value, time).Value()

	return &v
}

func (c *cache[K, V]) Unmarshal(key K, to interface{}) error {
	err := json.Unmarshal([]byte(c.GetString(key)), &to)
	return err
}

func Init() {
	GlobalCache = New[string, interface{}]()
}

func New[K comparable, V any]() Cache[K, V] {
	instance := ttlcache.New[K, V](
		ttlcache.WithTTL[K, V](15*time.Minute),
		ttlcache.WithDisableTouchOnHit[K, V](),
	)

	go instance.Start() // starts automatic expired item deletion

	return &cache[K, V]{cache: instance}
}

func interfaceToString(from interface{}) string {
	record, _ := json.Marshal(from)
	return string(record)
}
