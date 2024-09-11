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
var flagNtfy string
var flagNotification string
var Version string
var flagSysexOnly bool

func init() {
	flag.BoolVar(&flagUseSerial, "serial", false, "use serial")
	flag.BoolVar(&flagDebug, "debug", false, "debug")
	flag.BoolVar(&flagVersion, "version", false, "version")
	flag.BoolVar(&flagSysexOnly, "sysex", false, "sysex only")
	flag.StringVar(&flagMidiName, "midi", "", "midi name")
	flag.StringVar(&flagNtfy, "ntfy", "", "ntfy")
	flag.StringVar(&flagNotification, "on", "", "on")
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
		midicom.Run(flagMidiName, flagNtfy, flagNotification, flagSysexOnly)
	} else if flagUseSerial {
		serialcom.Run()
	}
}
