package main

import "time"

type MidiCommand struct {
	Changes map[string]uint8 `json:"changes"`
	Timeout float64          `json:"timeout"`
}

func (c MidiCommand) timeoutDuration() time.Duration {
	return time.Duration(c.Timeout * 1000 * 1000 * 1000)
}

type MidiCommands []MidiCommand
