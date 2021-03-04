//MIT License
//
//Copyright (c) 2020 Wout Slakhorst
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type RFIDReader struct {
	Codes    chan string
	Mappings Config
}

func (rfid RFIDReader) start() error {
	options := serial.OpenOptions{
		PortName:        "/dev/ttyS0",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 14,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		return err
	}

	go func() {
		// Make sure to close it later.
		defer port.Close()

		buf := make([]byte, 14)
		dBuf := make([]byte, 28)

		for {
			_, err := port.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			// create a double buf because the reading can be partial
			copy(dBuf[0:14], dBuf[14:])
			copy(dBuf[14:], buf[:])

			n, err := parseRFIDBytes(dBuf)
			if err != nil {
				log.Println(err)
				continue
			}
			tagId := fmt.Sprintf("%d", n)
			fmt.Printf("Scanned: %s\n", tagId)

			urlOrCommand, ok := rfid.Mappings.Mappings[tagId]
			if ok {
				fmt.Printf("Found track: %s\n", urlOrCommand)
				rfid.Codes <- urlOrCommand
				if urlOrCommand == "EXIT" {
					break
				}
				// when success, pause for a second
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return nil
}

func isHeadByte(s string) bool {
	return "\x02" == s
}

func isTailByte(s string) bool {
	return "\x03" == s
}

func checksum(cs []byte, d []byte) bool {
	stepSize := len(cs)
	current := uint32(0)
	for i := 0; i < len(d); i += stepSize {
		db, _ := strconv.ParseUint(string(d[i:i+2]), 16, 32)
		current = current ^ uint32(db)
	}
	n, _ := strconv.ParseUint(string(cs), 16, 32)

	if current != uint32(n) {
		log.Printf("%d vs %d", uint32(n), current)
		return false
	}
	return true
}

func parseRFIDBytes(doubleBuf []byte) (uint32, error) {

	for i := 0; i < 14; i++ {

		var buf = make([]byte, 14)
		copy(buf, doubleBuf[i:i+14])

		if !isHeadByte(string(buf[0])) {
			continue
			//return 0, fmt.Errorf("first byte was not 02 but %s\n", string(buf[0]))
		}

		if !isTailByte(string(buf[13])) {
			continue
			//return 0, fmt.Errorf("Last byte was not 03 but %s\n", string(buf[13]))
		}

		csBuf := buf[11:13]
		dBuf := buf[1:11]

		if !checksum(csBuf, dBuf) {
			continue
			//return 0, errors.New("Checksum failed")
		}

		s := string(buf[3:11])
		n, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			continue
			//return 0, fmt.Errorf("Error parsing %s: %v\n", s, err)
		}
		return uint32(n), nil
	}

	return 0, errors.New("no card found")
}
