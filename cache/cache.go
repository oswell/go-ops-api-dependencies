// Base cache
package cache

import (
	"fmt"

	"github.com/oswell/go-ops-api-dependencies/util"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/op/go-logging"
)

// Logger
var logger = logging.MustGetLogger("cache")

func init() {
	util.StderrFormatter("")
}

type Cache struct {
	Client *memcache.Client

	Servers    string
	Prefix     string
	Expiration int32
}

// cache.New("127.0.0.1:11211")
func New(servers string, expiration int32, prefix string) *Cache {
	c := Cache{Servers: servers, Expiration: expiration, Prefix: prefix}
	c.Client = memcache.New(c.Servers)
	return &c
}

func (c *Cache) Get(key string) ([]byte, error) {
	key_name := c.cache_key(key)
	logger.Infof("Getting cache entry for %s", key_name)

	item, err := c.Client.Get(key_name)
	if err != nil {
		logger.Debugf("No entry found, or entry expired for %s", key_name)
		return nil, err
	}

	return item.Value, nil
}

func (c *Cache) Set(key string, value []byte) {
	key_name := c.cache_key(key)
	logger.Infof("Setting new cache entry for %s (expiration in %d seconds)", key_name, c.Expiration)

	c.Client.Set(&memcache.Item{Key: key_name, Value: value, Expiration: c.Expiration})
}

func (c *Cache) Delete(key string) (error) {
	key_name := c.cache_key(key)
	logger.Infof("Deleting cache entry for %s", key_name)
	err := c.Client.Delete(key_name)

	if err != nil {
		logger.Errorf("Error clearing cache entry %s, %v", key_name, err)
		return err
	}

	return nil
}

func (c *Cache) cache_key(key string) (string) {
	return fmt.Sprintf("%s%s", c.Prefix, key)
}
