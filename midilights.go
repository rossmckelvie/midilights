package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosuri/uilive"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

var (
	channels map[string]LightChannel
	wg       sync.WaitGroup
)

func test() {
	var closeFn func()

	// Connect to light channels
	if runtime.GOARCH != "arm" {
		channels, closeFn = FactoryTestLightChannel()
	} else {
		channels, closeFn = FactoryPiLightChannel()
	}
	defer closeFn()

	// Flash all channels
	for _, channel := range channels {
		channel.on()
		time.Sleep(100 * time.Millisecond)
		channel.off()
	}

	// Play Music
	//go play("test.mp3")

	// Open test configuration & execute
	var cmds MidiCommands
	jsonFile, err := os.Open("test.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &cmds)

	writer := uilive.New()
	writer.Start()

	// Execute loaded configuration
	for _, command := range cmds {
		time.Sleep(command.timeoutDuration())
		wg.Add(2)
		go executeCommands(command)
		printChannelStatus(writer)
	}
	wg.Wait()
	writer.Stop()
}

func printChannelStatus(writer io.Writer) {
	output := ""

	keys := make([]string, 0)
	for k := range channels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		//var channelStatus = "OFF"
		//if channels[k].getStatus() == true {
		//	channelStatus = "ON"
		//}
		output += fmt.Sprintf("%s: %v\n", k, channels[k].getStatus())
	}

	fmt.Fprintf(writer, output)

	wg.Done()
}

func executeCommands(command MidiCommand) {
	for channelId, change := range command.Changes {
		if value, ok := channels[channelId]; ok {
			if change > 0 {
				go value.on()
			} else {
				go value.off()
			}
		}
	}
	wg.Done()
}
