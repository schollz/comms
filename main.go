package main

import (
	"flag"
	"fmt"

	"github.com/schollz/comms/src/midicom"
	"github.com/schollz/comms/src/serialcom"
	log "github.com/schollz/logger"
)

var flagUseSerial bool
var flagDebug bool
var flagMidiName string
var flagVersion bool
var Version string

func init() {
	flag.BoolVar(&flagUseSerial, "serial", false, "use serial")
	flag.BoolVar(&flagDebug, "debug", false, "debug")
	flag.BoolVar(&flagVersion, "version", false, "version")
	flag.StringVar(&flagMidiName, "midi", "", "midi name")
}

func main() {
	flag.Parse()
	if flagVersion {
		fmt.Printf("comms version %s\n", Version)
		return
	}
	if flagDebug {
		log.SetLevel("debug")
	} else {
		log.SetLevel("info")
	}

	if flagMidiName != "" {
		midicom.Run(flagMidiName)
	} else if flagUseSerial {
		serialcom.Run()
	}
}
