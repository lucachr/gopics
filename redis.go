/*
Redis utility functions for GoPics.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	redisMaxIdle     = 3
	redisDefaultAddr = ":6379"
)

var redisIdleTimeout = 240 * time.Second

// newPool creates a new connections pool for concurrent access
// to Redis.
func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// redisFlat flatens a struct for Redis HMSET
func redisFlat(key string, value interface{}) redis.Args {
	return redis.Args{}.Add(key).AddFlat(value)
}

// redisGetUser search for an user with the given username,
// it returns the user data if the user is found, otherwise,
// it returns a redis.ErrNil.
func redisGetUser(conn redis.Conn, username string) (*User, error) {
	val, err := redis.Values(conn.Do("HGETALL", userTag+username))
	switch {
	case err != nil:
		return nil, err
	case len(val) == 0:
		return nil, redis.ErrNil
	}

	usr := new(User)
	err = redis.ScanStruct(val, usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// redisGetPosts return a list of the latest one hundred posts in postSet.
// Posts are sorted by publishing date, starting from the latest one.
func redisGetPosts(conn redis.Conn, postSet string) ([]Post, error) {
	postNames, err := redis.Strings(conn.Do("ZREVRANGE", postSet, 0, 100))
	if err != nil {
		return nil, err
	}

	posts := []Post{}

	for _, name := range postNames {
		val, err := redis.Values(conn.Do("HGETALL", postTag+name))
		if err != nil {
			return nil, err
		}

		p := new(Post)
		err = redis.ScanStruct(val, p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *p)
	}

	return posts, nil
}
