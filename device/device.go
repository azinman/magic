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

package device

import (
	"net"

	localnet "github.com/azinman/magic/net"
)

type (
	// Kind represents an enum of device types
	Kind int

	// Device is a baseline interface for interacting with Magic-compatible
	// devices
	Device interface {
		Name() string
		UUID() string
		Kind() Kind
	}

	// Powerable describes devices that can be turned on or off
	Powerable interface {
		TurnOn() error
		TurnOff() error
	}

	// Networked applies to any LAN-based device that has an IP
	LANNetworked interface {
		// IPs (IPv4 and/or IPv6).
		IPs() []net.IP
		// MAC addresses should be available to all devices. Some devices
		// may change MAC addresses for privacy-protection, so this is
		// an array.
		MACs() []string
	}
)

const (
	KindTV Kind = iota
	KindPhone
	KindComputer
	KindUnknown
)

var registry = make(map[string]Device)
var ipCache = make(map[string]Device)
var macCache = make(map[string]Device)

// Register takes a device and retains a reference to it by its UUID.
func Register(device Device) {
	registry[device.UUID()] = device
	if networked, ok := device.(LANNetworked); ok {
		for _, ip := range networked.IPs() {
			ipCache[ip.String()] = device
		}
		found := false
		for _, mac := range networked.MACs() {
			macCache[mac] = device
			found = true
		}
		if !found {
			for _, ip := range networked.IPs() {
				if mac := localnet.LookupMACByIP(ip); mac != "" {
					macCache[mac] = device
					break
				}
			}
		}
	}
}

// FindDeviceByUUID returns any registered device for a given UUID.
func FindDeviceByUUID(uuid string) Device {
	return registry[uuid]
}

// FindDeviceByIP returns any registered device for a given IP.
func FindDeviceByIP(ip string) Device {
	return ipCache[ip]
}

// FindDeviceByMAC returns any registered device for a given ethernet MAC address.
func FindDeviceByMAC(mac string) Device {
	return macCache[mac]
}
