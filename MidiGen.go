package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gitlab.com/gomidi/midi/writer"
)

type printer struct{}

func main() {
	rand.Seed(time.Now().UnixNano())
	var (
		input    string
		numNotes int
		Rerr     error
	)
	//number of notes to generate
	fmt.Printf("How many notes would you like to generate?\n")
	fmt.Scanln(&input)
	if input == "0" || input == "" {
		numNotes = rand.Intn(10)
	} else {
		numNotes, Rerr = strconv.Atoi(input)
	}

	//check for /midi dir in working dir to save .mid to
	dir, Derr := os.Getwd()
	folderInfo, Derr := os.Stat(dir + "/midi")
	if Derr == nil {
		dir = filepath.Join(dir, folderInfo.Name())
	}
	if os.IsNotExist(Derr) {
		os.Mkdir("midi", 0755)
		Derr = nil
	}

	// file generation
	name := "randMidi_" + strconv.Itoa(time.Now().Hour()) + strconv.Itoa(time.Now().Minute()) + strconv.Itoa(time.Now().Second()) + ".mid"
	f := filepath.Join(dir, name)

	err := writer.WriteSMF(f, 1, func(wr *writer.SMF) error {
		wr.SetChannel(0)
		for n := 0; n < numNotes; n++ {
			var note uint8
			noteTemp := rand.Uint32()
			temp := make([]byte, 4)
			binary.LittleEndian.PutUint32(temp, noteTemp)

			note = temp[3]
			//TODO Clamp values to normal range
			//add key filters
			//make it play notes

			writer.NoteOn(wr, note, 50)
			wr.SetDelta(120)
			writer.NoteOff(wr, note)
		}
		writer.EndOfTrack(wr)

		return nil
	})

	if err != nil {
		fmt.Printf("could not write SMF file %v\n", f)
		return
	}

	if Rerr != nil {
		fmt.Printf("could not read imput \"%v\"\n", input)
	}

	if Derr != nil {
		fmt.Printf("could not find directory")
	}
}
