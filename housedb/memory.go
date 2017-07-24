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
	"sync"
)

type (
	memoryDb struct {
		data     map[string]string
		lock     *sync.Mutex
		watchers map[string]func(key string, value string)
	}
)

func newMemoryDb() *memoryDb {
	mdb := &memoryDb{
		data:     make(map[string]string),
		lock:     &sync.Mutex{},
		watchers: make(map[string]func(string, string)),
	}
	mdb.Watch("", func(key string, value string) {
		onChange(key, value)
	})
	return mdb
}

func (db *memoryDb) Get(key string) (string, error) {
	db.lock.Lock()
	defer db.lock.Unlock()
	return db.data[key], nil
}

func (db *memoryDb) Put(key string, value string) error {
	db.lock.Lock()
	db.data[key] = value
	db.lock.Unlock()
	// Notify off-lock
	for callbackKey, callback := range db.watchers {
		if callbackKey == key || callbackKey == "" {
			callback(key, value)
		}
	}
	return nil
}

// Watch calls the callback with new values for a given key. If key is
// empty then all values are turned.
func (db *memoryDb) Watch(key string, callback func(key string, value string)) error {
	// Ghetto watcher for now
	db.watchers[key] = callback
	return nil
}

func (db *memoryDb) Close() error {
	return nil
}
