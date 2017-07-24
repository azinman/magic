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

package wifi

import (
	"fmt"
	"log"

	"github.com/azinman/magic/context"
	"github.com/google/gopacket/pcap"
)

type (
	// Wifi Low Level Event
	WifiEvent struct {
		name string
	}
)

func (e *WifiEvent) Name() string {
	return e.name
}
func (e *WifiEvent) Data() interface{} {
	return nil
}

// MonitorWifi emits events
func MonitorWifi() <-chan context.LowLevelEvent {
	fmt.Println("Making stuff")
	out := make(chan context.LowLevelEvent)
	go func() {
		findDevice()
		close(out)
	}()
	return out
}

func findDevice() {
	fmt.Println("Findin all devices")
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Print device information
	fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
}
