package domain

import "github.com/go-redis/redis"

type TermsCache struct {
	redis *redis.Client
}

func NewTermsCache() TermsCache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return TermsCache{
		redis: client,
	}
}

func (c *TermsCache) IncreaseTermCount(term string) {
	cmd := c.redis.Incr(term)
	err := cmd.Err()

	if err != nil {
		panic(err)
	}
}
