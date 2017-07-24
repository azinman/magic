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

package context

import (
	"net"

	"github.com/azinman/magic/device"
)

type (
	dnsLookup struct {
		host   string
		device *device.Device
	}

	unknownNetworkDevice struct {
		ips []net.IP
		mac string
	}
)

func (d *unknownNetworkDevice) IPs() []net.IP {
	return d.ips
}

func (d *unknownNetworkDevice) MAC() string {
	if d.mac != "" {
		return d.mac
	}
	// TODO: Lookup in cache
	return ""
}
