package midicom

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	log "github.com/schollz/logger"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

var isConnected = false

func isAvailable(filterMidiName string) bool {
	return strings.Contains(strings.ToLower(midi.GetInPorts().String()), strings.ToLower(filterMidiName))
}

func doConnection(filterMidiName string) (stop func(), err error) {
	var midiInput drivers.In
	ins := midi.GetInPorts()
	if len(ins) == 0 {
		log.Error("no input devices")
		return
	}

	for _, in := range ins {
		log.Tracef("found input: '%s'", in.String())
		if strings.Contains(strings.ToLower(in.String()), strings.ToLower(filterMidiName)) {
			midiInput = in
			break
		}
	}
	if midiInput == nil {
		log.Error("no input devices")
		return
	}

	// listen to midi
	stop, err = midi.ListenTo(midiInput, func(msg midi.Message, timestamps int32) {
		var bt []byte
		var ch, key, vel uint8
		switch {
		case msg.GetSysEx(&bt):
			fmt.Printf("%s", bt)
			if ntfyMessage != "" && ntfyTopic != "" && strings.Contains(string(bt), ntfyMessage) {
				log.Infof("sending ntfy to %s", ntfyTopic)
				http.Post("https://ntfy.sh/"+ntfyTopic, "text/plain",
					strings.NewReader("comms: "+string(bt)))
			}
		case msg.GetNoteStart(&ch, &key, &vel) && !sysexOnly:
			log.Infof("note_on=%s, ch=%v, vel=%v\n", midi.Note(key), ch, vel)
		case msg.GetNoteEnd(&ch, &key) && !sysexOnly:
			log.Infof("note_off=%s, ch=%v\n", midi.Note(key), ch)
		default:
			// ignore
		}
	}, midi.UseSysEx(), midi.SysExBufferSize(4096))

	if err != nil {
		log.Error(err)
		return
	}

	isConnected = true
	log.Infof("connected to\n\t'%s'", midiInput.String())
	return
}

var filterMidiName string
var ntfyTopic string
var ntfyMessage string
var sysexOnly bool

func Run(name string, ntfy string, notification string, sysex bool) {
	filterMidiName = name
	ntfyTopic = ntfy
	ntfyMessage = notification
	sysexOnly = sysex

	var err error
	done := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			done <- struct{}{}
			time.Sleep(100 * time.Millisecond)
			log.Infof("captured %v, exiting..", sig)
			os.Exit(0)
		}
	}()

	var stopFunc func()
	go func() {
		// check if midi is connected
		for {
			if isConnected {
				if !isAvailable(filterMidiName) {
					isConnected = false
					midi.CloseDriver()
				}
			} else {
				if isAvailable(filterMidiName) {
					stopFunc, err = doConnection(filterMidiName)
					if err != nil {
						log.Error(err)
					}
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		select {
		case <-done:
			stopFunc()
		}
	}

}
