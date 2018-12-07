package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"log"
	"os"
)

type LightChannel interface {
	setChannelId(id string)
	getChannelId() string

	setPinId(id uint8)
	getPinId() uint8

	getStatus() bool

	on()
	off()
}

type TestLightChannel struct {
	channelId string
	pinId     uint8
	onStatus  bool
}
func (t TestLightChannel) setChannelId(id string) { t.channelId = id }
func (t TestLightChannel) getChannelId() string   { return t.channelId }
func (t TestLightChannel) setPinId(id uint8)      { t.pinId = id }
func (t TestLightChannel) getPinId() uint8        { return t.pinId }
func (t TestLightChannel) on()                    { t.onStatus = true }
func (t TestLightChannel) off()                   { t.onStatus = false }
func (t TestLightChannel) logPrefix() string      { return fmt.Sprintf("[%s P%d]", t.channelId, t.pinId) }
func (t TestLightChannel) getStatus() bool        { return t.onStatus }
func NewTestLightChannel(channelId string, pinId uint8) TestLightChannel {
	channel := TestLightChannel{channelId: channelId, pinId: pinId, onStatus: false}
	return channel
}
func FactoryTestLightChannel() (map[string]LightChannel, func()) {
	channels := make(map[string]LightChannel)

	channels["right-tree"] = NewTestLightChannel("right-tree", 0)
	channels["left-tree"] = NewTestLightChannel("left-tree", 1)
	channels["front-door"] = NewTestLightChannel("front-door", 2)
	channels["laser"] = NewTestLightChannel("laser", 3)
	channels["garage"] = NewTestLightChannel("garage", 4)
	channels["left-window"] = NewTestLightChannel("left-window", 5)
	channels["right-window"] = NewTestLightChannel("right-window", 6)
	channels["icicles"] = NewTestLightChannel("icicles", 7)

	return channels, func() { log.Println("TestFactory Closing Down") }
}

var gpio2bcm = map[uint8]uint8{
	0: 17,
	1: 18,
	2: 27,
	3: 22,
	4: 23,
	5: 24,
	6: 25,
	7: 4,
}

type PiLightChannel struct {
	channelId string
	pinId     uint8
	onStatus  bool
	pin       rpio.Pin
}

func (p PiLightChannel) setChannelId(id string) { p.channelId = id }
func (p PiLightChannel) getChannelId() string   { return p.channelId }
func (p PiLightChannel) setPinId(id uint8)      { p.pinId = id }
func (p PiLightChannel) getPinId() uint8        { return p.pinId }
func (p PiLightChannel) on()                    { p.onStatus = true; p.pin.High() }
func (p PiLightChannel) off()                   { p.onStatus = false; p.pin.Low() }
func (p PiLightChannel) logPrefix() string      { return fmt.Sprintf("[%s P%d]", p.channelId, p.pinId) }
func (p PiLightChannel) getStatus() bool        { return p.onStatus }
func NewPiLightChannel(channelId string, pinId uint8) PiLightChannel {
	channel := PiLightChannel{channelId: channelId, pinId: pinId, onStatus: false}

	pin := rpio.Pin(gpio2bcm[pinId])
	pin.Output()
	channel.pin = pin

	return channel
}

func FactoryPiLightChannel() (map[string]LightChannel, func()) {
	if os.Geteuid() != 0 {
		log.Fatal("This program have to be run as root, or SUID/GUID set to 0 on execution!")
		os.Exit(1)
	}

	closeFn := InitGPIO()

	channels := make(map[string]LightChannel)

	channels["right-tree"] = NewPiLightChannel("right-tree", 0)
	channels["left-tree"] = NewPiLightChannel("left-tree", 1)
	channels["front-door"] = NewPiLightChannel("front-door", 2)
	channels["laser"] = NewPiLightChannel("laser", 3)
	channels["garage"] = NewPiLightChannel("garage", 4)
	channels["left-window"] = NewPiLightChannel("left-window", 5)
	channels["right-window"] = NewPiLightChannel("right-window", 6)
	channels["icicles"] = NewPiLightChannel("icicles", 7)

	return channels, closeFn
}

func InitGPIO() (disconnect func()) {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return func() {
		rpio.Close()
	}
}
