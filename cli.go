package main

import (
	"fmt"
	"runtime"
	"time"

	docopt "github.com/docopt/docopt-go"
)

func main() {
	usage := `Control your midi lights from the command line.

Usage:
	ml configure
	ml lights (on|off|flash) [--channel=CHANNEL_ID]
  
Options:
	-h, --help  Show this screen.`

	arguments, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}

	// Connect to light channels
	var channels map[string]LightChannel
	var closeFn func()
	if runtime.GOARCH != "arm" {
		channels, closeFn = FactoryTestLightChannel()
	} else {
		channels, closeFn = FactoryPiLightChannel()
	}
	defer closeFn()

	if arguments["lights"] == true {
		err = lights(arguments, channels)
	}

	if err != nil {
		fmt.Println("Something went wrong.", err)
	}
}

func lights(arguments docopt.Opts, channels map[string]LightChannel) error {
	
	// Filter for channel
	if arguments["--channel"] != nil {
		var channelSpecific = arguments["--channel"].(string)
		channels = map[string]LightChannel{channelSpecific: channels[channelSpecific]}
	}

	// Toggle perform action
	for _, channel := range channels {
		fmt.Println(channel.getChannelId())
		channel.on()
		if arguments["on"] == true {
			channel.on()
		} else if arguments["off"] == true {
			channel.off()
		} else if arguments["flash"] == true {
			for i := 0; i < 3; i++ {
				channel.on()
				time.Sleep(200 * time.Millisecond)
				channel.off()
				time.Sleep(200 * time.Millisecond)
			} 
		}
	}

	return nil
}