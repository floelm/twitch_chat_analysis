package pkg

import (
	"github.com/go-redis/redis"
	"strings"
)

type EmoticonsCache struct {
	redis *redis.Client
}

func NewEmoticonsCache() EmoticonsCache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return EmoticonsCache{
		redis: client,
	}
}

func (c *EmoticonsCache) Exists(name string) bool {
	cmd := c.redis.Exists(strings.ToLower(name))
	exists, _ := cmd.Result()

	return exists == 1
}

func (c *EmoticonsCache) Increase(name string) {
	c.redis.Incr(strings.ToLower(name))
}
