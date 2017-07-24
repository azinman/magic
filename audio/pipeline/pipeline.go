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

import "log"

type PipeType int

const (
	Source PipeType = iota
	Sink
	Transformer
)

func (s PipeType) String() string {
	switch s {
	case Source:
		return "source"
	case Sink:
		return "sink"
	case Transformer:
		return "transformer"
	}
	return "unknown"
}

type Encoding struct {
	Format     string
	Channels   int
	SampleRate int
}

// Pipe is an individual segment in a pipeline. They read from the InputPipe's
// output and write to its feedback -- this happens on connecting the input.
// Pipes write to their output on OutputChan, and receive feedback from their
// output on Feedback chan.
type Pipe struct {
	Type           PipeType
	InputEncoding  Encoding
	OutputEncoding Encoding
	Description    string
	InputPipe      *Pipe
	FeedbackChan   <-chan Packet
	OutputChan     <-chan Packet
}

func (p *Pipe) Init(pipeType PipeType, inputEnc Encoding, outputEnc Encoding, desc string) {
	p.Type = pipeType
	p.InputEncoding = inputEnc
	p.OutputEncoding = outputEnc
	p.Description = desc
	p.FeedbackChan = make(<-chan Packet)
	p.OutputChan = make(<-chan Packet)
}

func (p *Pipe) SetInput(input *Pipe) {
	p.InputPipe = input

	// Input data
	go func() {
		for packet := range p.InputPipe.OutputChan {
			p.HandleInput(packet)
		}
		// close(p.OutputChan)
	}()

	// Feedback chan
	go func() {
		for packet := range p.FeedbackChan {
			newFeedback := p.HandleFeedback(packet)
			if newFeedback != nil {
				p.InputPipe.FeedbackChan <- newFeedback
			}
		}
		// close(p.InputPipe.FeedbackChan)
	}()
}

func (p *Pipe) HandleInput(packet Packet) *Packet {
	log.Printf("handleInput %v", packet)
	return &packet
}

func (p *Pipe) HandleFeedback(packet Packet) *Packet {
	log.Printf("handleFeedback %v", packet)
	return nil
}
