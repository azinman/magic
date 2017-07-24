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

package pipeline

import "fmt"

type PacketType int

const (
	QueuePlayAt PacketType = iota
	DidPlayAt
	WillPlayAt
	Dropped
)

func (s PacketType) String() string {
	switch s {
	case QueuePlayAt:
		return "queuePlayAt"
	case DidPlayAt:
		return "didPlayAt"
	case WillPlayAt:
		return "willPlayAt"
	case Dropped:
		return "dropped"
	}
	return "unknown"
}

type Packet struct {
	PacketNum int
	Type      PacketType
}

type QueuePacket struct {
	Packet
	MediaTime int64
	Data      *[]byte
}

func NewQueuePacket(num int, mediaTime int64, data *[]byte) QueuePacket {
	p := &QueuePacket{
		MediaTime: mediaTime,
		Data:      data,
	}
	p.PacketNum = num
	p.Type = QueuePlayAt
	return p
}

func (p QueuePacket) String() string {
	return fmt.Sprintf("[QueuePacket #%v time=%v data=%v bytes]",
		p.PacketNum, p.MediaTime, len(*p.Data))
}

type DidPlayPacket struct {
	Packet
	PlayedMediaTime    int64
	RequestedMediaTime int64
}

func (p DidPlayPacket) String() string {
	return fmt.Sprintf("[DidPlayPacket #%v playedTime=%v requestedTime=%v diff=%v]",
		p.PacketNum, p.PlayedMediaTime, p.RequestedMediaTime,
		p.PlayedMediaTime-p.RequestedMediaTime)
}

type WillPlayPacket struct {
	Packet
	MediaTime          int64
	RequestedMediaTime int64
}

func (p WillPlayPacket) String() string {
	return fmt.Sprintf("[WillPlayPacket #%v time=%v requestedTime=%v diff=%v]",
		p.PacketNum, p.MediaTime, p.RequestedMediaTime,
		p.MediaTime-p.RequestedMediaTime)
}

type DroppedPacket struct {
	Packet
	RequestedMediaTime int64
}

func (p DroppedPacket) String() string {
	return fmt.Sprintf("[DroppedPacket #%v requestedTime=%v]",
		p.PacketNum, p.RequestedMediaTime)
}
