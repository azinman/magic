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

package sources

import (
	"fmt"
	"log"
	"os"

	"github.com/azinman/magic/audio/pipeline"
	"io"
)

const (
	SampleRate = 44100
	ReadBufferLen = SampleRate
)

type NamedPipeSource struct {
	pipeline.Pipe
	File *os.File
}

func NewNamedPipeSource(pipePath string, encoding pipeline.Encoding) (*NamedPipeSource, error) {
	log.Println("Opening pipe ", pipePath)
	f, err := os.Open(pipePath)
	if err != nil {
		return nil, err
	}
	pipe := &NamedPipeSource{
		File: f,
	}
	pipe.Init(pipeline.Source, encoding, encoding, fmt.Sprintf("Named pipe source from %v", pipePath))
	go func(f *os.File, outputChan <-chan pipeline.Packet) {
		num := 0
		mediaTime := 0
		for {
			b := make([]byte, ReadBufferLen)
			n, err := f.Read(b)
			switch err {
			case nil:
				data := b
				if n < ReadBufferLen {
					data = b[:n]
				}
				num++
				pipeline.NewQueuePacket(num, mediaTime, data)
				mediaTime +=
				continue
			case io.EOF:
				log.Print("EOF:", pipe.Description)
				break
			default:
				panic(err)
			}

		}

	}()

	return pipe, nil
}
