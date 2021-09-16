package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gitlab.com/gomidi/midi/writer"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var (
		noteIn   string
		numNotes int
		keyIn    string
		key      int
		octIn    string
		oct      int
		noteRerr error
		octRerr  error
	)

	//number of notes to generate
	fmt.Printf("How many notes would you like to generate?\n")
	fmt.Scanln(&noteIn)
	if noteIn == "0" || noteIn == "" {
		//random number of notes
		numNotes = rand.Intn(64)
	} else {
		numNotes, noteRerr = strconv.Atoi(noteIn)
	}
	if noteRerr != nil {
		fmt.Printf("could not read input %v\n", noteIn)
		os.Exit(1)
	}

	//select major or minor
	fmt.Printf("major or minor? (0 for major, 1 for minor) default is major\n")
	fmt.Scanln(&keyIn)
	if keyIn == "1" || keyIn == "minor" || keyIn == "min" {
		key = 1
	} else {
		key = 0
	}

	//octaves to span (default is 1)
	fmt.Printf("how many octaves would you like to span?\n")
	fmt.Scanln(&octIn)
	oct, octRerr = strconv.Atoi(octIn)
	if oct == 0 {
		octRerr = errors.New("Value: value cannot be 0")
	}
	if octRerr != nil {
		fmt.Printf("could not read input %v\n", octIn)
		os.Exit(2)
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

	//write midi file
	err := writer.WriteSMF(f, 1, func(wr *writer.SMF) error {
		wr.SetChannel(0)
		for n := 0; n < numNotes; n++ {
			//TODO
			//make it play notes
			//randomise SetDeltas for different length of notes

			note := note(key, oct)

			writer.NoteOn(wr, note, 50)
			//length of note on
			wr.SetDelta(240)

			writer.NoteOff(wr, note)
			//length of note off
			wr.SetDelta(120)
		}
		writer.EndOfTrack(wr)

		return nil
	})

	if err != nil {
		fmt.Printf("could not write SMF file %v\n", f)
		return
	}

	if Derr != nil {
		fmt.Printf("could not find directory \n")
		return
	}
}

func note(m int, octaves int) uint8 {
	var (
		n uint8
		o int
	)
	oct := octaves
	r := rand.Intn(8)
	//random number of loops to add octave
	if octaves > 1 {
		o = rand.Intn(oct) + 1
	} else {
		o = octaves
	}
	if m == 0 {
		switch r {
		//C
		case 0:
			n = 72
		//D
		case 1:
			n = 74
		//E
		case 2:
			n = 76
		//F
		case 3:
			n = 77
		//G
		case 4:
			n = 79
		//A
		case 5:
			n = 81
		//B
		case 6:
			n = 83
		//C
		case 7:
			n = 84
		}
	} else {
		switch r {
		//C
		case 0:
			n = 72
		//D
		case 1:
			n = 74
		//Eb
		case 2:
			n = 75
		//F
		case 3:
			n = 77
		//G
		case 4:
			n = 78
		//Ab
		case 5:
			n = 81
		//Bb
		case 6:
			n = 82
		//C
		case 7:
			n = 84
		}
	}
	for i := 1; i < o; i++ {
		n += 12
	}
	return n
}
