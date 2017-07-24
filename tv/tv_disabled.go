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

package tv

import (
	"fmt"

	// "github.com/chbmuc/cec"
	"github.com/azinman/magic/device"
)

// TV represents your TV
type TV struct {
	// c *cec.Connection
}

// NewTV opens a connection to the TV and returns a struct representing i/
func NewTV() (*TV, error) {
	//c, err := cec.Open("", "cec.go")
	// if err != nil {
	// 	return nil, err
	// }
	// return &TV{c: c}, nil
	return nil, fmt.Errorf("Disabled until jail in freebsd")
}

// Info tells things about the TV
func (tv *TV) Info() {
	// log.Print("List: ", tv.c.List())
	// log.Print("Active Devices: ", tv.c.GetActiveDevices())
}

func (tv *TV) Name() string {
	return "TV"
}

func (tv *TV) UUID() string {
	return "0db0155a-8d06-40aa-8faf-23c7a4ac1347"
}

func (tv *TV) Kind() device.Kind {
	return device.KindTV
}

// TurnOn turns off the CEC-attached TV
func (tv *TV) TurnOn() error {
	// return tv.c.PowerOn(0)
	return fmt.Errorf("Disabled until jail in freebsd")
}

// TurnOff turns off the CEC-attached TV
func (tv *TV) TurnOff() error {
	// return tv.c.Standby(0)
	return fmt.Errorf("Disabled until jail in freebsd")
}
