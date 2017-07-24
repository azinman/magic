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

type (
	// Db represents the abstracted housedb, which might use multiple underlying
	// implementations. Currently it's redis.
	Db interface {
		Get(key string) (string, error)
		Put(key string, value string) error
		Watch(key string, callback func(key string, value string)) error
		Close() error
	}
)

// NewMemory returns an in-memory database, not persistent.
func NewMemory() (Db, error) {
	return newMemoryDb(), nil
}

// NewRedis creates a new housedb using the flag-passed info
func NewRedis() (Db, error) {
	return newRedisDBFromFlags()
}
