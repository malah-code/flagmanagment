package cache

import "github.com/redis/go-redis/v9"

// Client wraps a go-redis client.
type Client struct {
	rdb *redis.Client
}

// NewClient initializes a Redis client instance.
func NewClient(addr, password string, db int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Client{rdb: rdb}
}

// Raw returns the underlying go-redis Client.
func (c *Client) Raw() *redis.Client {
	return c.rdb
}
