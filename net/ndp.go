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

package net

import (
	"os/exec"
	"strings"
	"sync"
)

var ndp = &ndpCache{
	table: make(ndpTable),
}

type (
	ndpTable map[string]string

	ndpCache struct {
		sync.RWMutex
		table ndpTable

		UpdatedCount int
	}
)

func (c *ndpCache) Refresh() {
	c.Lock()
	defer c.Unlock()

	c.table = newTable()
	c.UpdatedCount++
}

func (c *ndpCache) Search(ip string) string {
	c.RLock()
	defer c.RUnlock()

	mac, ok := c.table[ip]

	if !ok {
		c.RUnlock()
		c.Refresh()
		c.RLock()
		mac = c.table[ip]
	}

	return mac
}

func newTable() ndpTable {
	data, err := exec.Command("ndp", "-an").Output()
	if err != nil {
		return nil
	}

	var table = make(ndpTable)
	for _, line := range strings.Split(string(data), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		ipv6 := fields[0]
		mac := fields[1]
		table[ipv6] = mac
	}
	return table
}
