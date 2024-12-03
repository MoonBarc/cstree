package main

import (
	"encoding/binary"
	"log"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	width  = 150
	height = 1
)

var leds *ws2811.WS2811

func SetupLEDStrip() {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = ws2811.DefaultBrightness
	opt.Channels[0].LedCount = width * height

	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		log.Fatalln("failed to get device", err)
	}
	err = dev.Init()
	if err != nil {
		log.Fatalln("failed to init device", err)
	}

	leds = dev
}

func SetLEDs(states []byte) error {
	ls := leds.Leds(0)
	for i := 0; i < len(states); i += 4 {
		px := []byte{states[i+1], states[i+2], states[i+3], 0x00}
		ls[i/4] = binary.NativeEndian.Uint32(px)
	}
	return leds.Render()
}
