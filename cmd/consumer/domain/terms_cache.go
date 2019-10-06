package domain

import (
	"github.com/go-redis/redis"
	"sort"
)

type kv struct {
	Key   string
	Value int
}

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

func (c *TermsCache) GetAll() []kv {
	cmd := c.redis.Keys("*") //TODO: this is generally a bad idea :)
	keysSlice, err := cmd.Result()
	if err != nil {
		panic(err)
	}

	resultMap := make(map[string]int, 0)

	for _, redisKey := range keysSlice {
		getResult := c.redis.Get(redisKey)
		count, err := getResult.Int()
		if err != nil {
			panic(err)
		}
		resultMap[redisKey] = count
	}

	var sorted []kv
	for k, v := range resultMap {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	return sorted
}
