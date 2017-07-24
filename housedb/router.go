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
	"log"
	"strings"

	"github.com/azinman/magic/device"
)

func onChange(key string, value string) {
	s := strings.Split(key, ".")
	if len(s) != 2 {
		log.Print("Didn't understand key ", key)
		return
	}
	uuid, key := s[0], s[1]
	d := device.FindDeviceByUUID(uuid)
	if d == nil {
		log.Print("Can't find device uuid ", uuid)
		return
	}
	switch key {
	case "power":
		if pd, ok := d.(device.Powerable); ok {
			switch value {
			case "on":
				pd.TurnOn()
			case "off":
				pd.TurnOff()
			default:
				log.Print("Didn't understand ", value)
			}
		} else {
			log.Print("Got a power event for a non-powerable device")
		}
	default:
		log.Print("Don't know what to do with ", key)
	}
}
