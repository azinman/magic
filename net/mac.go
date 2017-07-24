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
	"log"
	"net"

	"github.com/mostlygeek/arp"
)

func LookupMACByString(ip string) string {
	if lookup := arp.Search(ip); lookup != "" {
		return lookup
	}
	if lookup := ndp.Search(ip); lookup != "" {
		return lookup
	}
	parsedIP := net.ParseIP(ip)
	return LookupMACByIP(parsedIP)
}

func LookupMACByIP(ip net.IP) string {
	if arp.CacheUpdateCount() == 0 {
		arp.CacheUpdate()
	}
	str := ip.String()
	if lookup := arp.Search(str); lookup != "" {
		return lookup
	}
	if ndp.UpdatedCount == 0 {
		ndp.Refresh()
	}
	if lookup := ndp.Search(str); lookup != "" {
		return lookup
	}

	if len(ip) == net.IPv4len {
		arp.CacheUpdate()
		if lookup := arp.Search(str); lookup != "" {
			return lookup
		}
		log.Printf("Couldn't find MAC address for %s\n", str)
	}
	ndp.Refresh()
	if lookup := ndp.Search(str); lookup != "" {
		return lookup
	}
	log.Printf("Couldn't find MAC address for %s\n", str)
	return ""
}
