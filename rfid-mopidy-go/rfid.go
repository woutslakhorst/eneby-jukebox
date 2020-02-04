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
	"fmt"
	"log"

	"github.com/karalabe/hid"
)

type RFIDReader struct {
	Codes    chan string
	Mappings Mappings
}

func (rfid RFIDReader) start() error {
	// Iterate through available Devices, finding all that match a known VID/PID.
	devs := hid.Enumerate(uint16(0xffff), uint16(0x0035))

	dev, err := devs[0].Open()
	if err != nil {
		return err
	}

	go func() {
		// close
		defer dev.Close()

		buf := make([]byte, 3)

		for {
			total := make([]byte, 0)
			for j := 0; j < 22; j++ {
				readBytes, err := dev.Read(buf)
				if err != nil {
					log.Fatalf("Read returned an error: %s", err.Error())
				}
				if readBytes == 0 {
					log.Fatalf("HID device returned 0 bytes of data.")
				}
				total = append(total, buf...)
			}

			// put numbers together to 10 chars
			number := make([]byte, 0)
			for _, n := range total {
				if n != 0x0 {
					switch n {
					case 30: number = append(number, '1')
					case 31: number = append(number, '2')
					case 32: number = append(number, '3')
					case 33: number = append(number, '4')
					case 34: number = append(number, '5')
					case 35: number = append(number, '6')
					case 36: number = append(number, '7')
					case 37: number = append(number, '8')
					case 38: number = append(number, '9')
					case 39: number = append(number, '0')
					case 40: //ignore newline
					}
				}
			}

			// keycodes to ascii
			fmt.Printf("Scanned: %s\n", string(number))

			urlOrCommand, ok := rfid.Mappings.Mappings[string(number)]
			if ok {
				fmt.Printf("Found track: %s\n", urlOrCommand)
				rfid.Codes <- urlOrCommand
				if urlOrCommand == "EXIT" {
					break
				}
			}

			// needed for reset?
			buf := make([]byte, 3)
			readBytes, err := dev.Write(buf)
			if err != nil {
				log.Fatalf("Write returned an error: %s", err.Error())
			}
			if readBytes == 0 {
				log.Fatalf("HID device received 0 bytes of data.")
			}
		}
	}()
	return nil
}