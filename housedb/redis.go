// Copyright Aaron Zinman 2017, 2018
// Copyright Duck Research LLC 2017, 2018
// All rights reserved.
//
// This file is part of Magichaus.
//
// Magichaus is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Magichaus is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Magichaus.  If not, see <http://www.gnu.org/licenses/>.

package housedb

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	redisURL      = flag.String("redisURL", "", "Redis URL")
	redisPassword = flag.String("redisPassword", "", "Redis password")
)

type (
	redisDb struct {
		pool     *redis.Pool
		watchers map[string]func(key string, value string)
	}
)

func newRedisDBFromFlags() (*redisDb, error) {
	if redisURL == nil {
		return nil, errors.New("-redisURL must be defined")
	}
	pw := ""
	if redisPassword != nil {
		pw = *redisPassword
	}
	return newRedisDB(*redisURL, pw), nil
}

// New creates a new housedb using the following redis connection info
func newRedisDB(redisURL string, redisPassword string) *redisDb {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			if err != nil {
				return nil, err
			}
			if redisPassword != "" {
				if _, err = c.Do("AUTH", redisPassword); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	c := pool.Get()
	defer c.Close()
	if c.Err() != nil {
		log.Fatal("Unable to connect to redis: ", c.Err())
		return nil
	}
	rdb := &redisDb{pool: pool, watchers: make(map[string]func(string, string))}
	rdb.Watch("", func(key string, value string) {
		onChange(key, value)
	})
	return rdb
}

func (db *redisDb) Get(key string) (string, error) {
	c := db.pool.Get()
	defer c.Close()
	return redis.String(c.Do("GET", key))
}

func (db *redisDb) Put(key string, value string) error {
	c := db.pool.Get()
	defer c.Close()
	_, err := c.Do("SET", key, value)
	if err != nil {
		return err
	}
	for callbackKey, callback := range db.watchers {
		if callbackKey == key || callbackKey == "" {
			callback(key, value)
		}
	}
	return nil
}

// Watch calls the callback with new values for a given key. If key is
// empty then all values are turned.
func (db *redisDb) Watch(key string, callback func(key string, value string)) error {
	// Ghetto watcher for now
	db.watchers[key] = callback
	return nil
}

func (db *redisDb) Close() error {
	return db.pool.Close()
}
