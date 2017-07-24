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

type (
	LowLevelEvent interface {
		Name() string
		Data() interface{}
	}
	HighLevelEvent interface {
		Name() string
		SourceEvents() []LowLevelEvent
		Data() interface{}
	}
)

/*
Streams:
	from traffic be able to tag as phone/mac/tv/etc
	cache owners of each device, manual for now.
	from unbound -> ip, cache IP lookup to provide linked context for packet filter
	from ip -> tag purpose of connection as catan
	from purpose, add to state of owner (person) of device
	from purpose, add to state of detected room


	each room gets a pi (zero?) running in wifi/bluetooth monitoring mode
		as it detects changings in RSSI, determine closest and add to a room


	can we flag data as interactive vs passive?
		record default traffic, ignore that
		everything else is active?
		{google, yelp, maps, catan} search -> active

	is there a network call every time i unlock phone? (prob not)
	is there a different network status (wifi active and not sleeping?)


Nene:
	automatically know when i go to kitchen shortly after timer goes off elsewhere, and turn off timer if one is ringing
	allow timers to fire beyond room they're in, if i move rooms
*/
