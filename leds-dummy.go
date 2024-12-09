//go:build !linux

package main

import (
	"log"
)

func SetupLEDStrip() {
	log.Println("warn: setting up dummy LED strip")
}

func SetLEDs(states []byte) error {
	// this is a dummy LED setting function
	return nil
}
