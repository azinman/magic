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

package sinks

import (
	"fmt"
	"log"
	"os"

	"github.com/azinman/magic/audio/pipeline"
)

type NamedPipeSink struct {
	pipeline.Pipe
	File *os.File
}

func NewNamedPipeSink(pipePath string, encoding pipeline.Encoding) (*NamedPipeSink, error) {
	log.Println("Opening pipe ", pipePath)
	f, err := os.Open(pipePath)
	if err != nil {
		return nil, err
	}
	pipe := &NamedPipeSink{
		File: f,
	}
	pipe.Init(pipeline.Sink, encoding, encoding, fmt.Sprintf("Named pipe sink to %v", pipePath))
	return pipe, nil
}
