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
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mappings   map[string]string `yaml:mappings`
	StartSound string            `yaml:startsound`
}

func main() {
	// config options
	flagSet := pflag.NewFlagSet("config", pflag.ContinueOnError)
	flagSet.String("mappings_file", "/etc/rfid-mopidy/rfid-mopidy.yaml", "Mopidy-rfid Config file")
	flagSet.String("mopidy_api", "http://localhost:6680", "Mopidy API endpoint")
	flagSet.Bool("stdin", false, "Read from stdin, use when rfid reader prints Codes as keyboard.")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		if err != pflag.ErrHelp {
			fmt.Printf("error reading args: %v", err)
			os.Exit(1)
		}
	}

	// read Config file
	fn, _ := flagSet.GetString("mappings_file")
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Printf("error reading %s: %s\n", fn, err.Error())
		os.Exit(1)
	}

	config := Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("error reading %s: %s\n", fn, err.Error())
	}

	fmt.Printf("startsound: %s\n", config.StartSound)
	for k, v := range config.Mappings {
		fmt.Printf("%s -> %s\n", k, v)
	}

	// needed channels
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	codes := make(chan string, 1)
	idle := make(chan bool, 1)

	// register interrupt
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// register halt command
	if err := monitorPin5(); err != nil {
		log.Fatalf("error starting pin monitoring: %s", err.Error())
	}

	// mopidy client
	mopidyClient := MopidyClient{
		RPCAddress: "http://localhost:6680/mopidy/rpc",
	}

	mopidyClient.waitForOk()
	mopidyClient.clearPlaylist()
	mopidyClient.setTracklist(fmt.Sprintf("file:%s", config.StartSound))
	mopidyClient.play()
	mopidyClient.detectIdle(idle)

	// rfid handler
	rfid := RFIDReader{
		Codes:    codes,
		Mappings: config,
	}
	if err := rfid.start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case code := <-codes:
				fmt.Printf("received code: %s\n", code)
				switch code {
				case "EXIT":
					done <- true
					break
				case "STOP":
					if err := mopidyClient.stop(); err != nil {
						log.Printf("mopidy.stop: %s\n", err.Error())
					}
				case "NEXT":
					if err := mopidyClient.next(); err != nil {
						log.Printf("mopidy.next: %s\n", err.Error())
					}
				case "PREV":
					if err := mopidyClient.previous(); err != nil {
						log.Printf("mopidy.previous: %s\n", err.Error())
					}
				default:
					// stop, clear, add, play
					if err := mopidyClient.stop(); err != nil {
						log.Printf("mopidy.stop: %s\n", err.Error())
					}
					if err := mopidyClient.clearPlaylist(); err != nil {
						log.Printf("mopidy.clear: %s\n", err.Error())
					}
					if err := mopidyClient.setTracklist(code); err != nil {
						log.Printf("mopidy.set: %s\n", err.Error())
					}
					if err := mopidyClient.shufflePlaylist(); err != nil {
						log.Printf("mopidy.shuffle: %s\n", err.Error())
					}
					if err := mopidyClient.play(); err != nil {
						log.Printf("mopidy.play: %s\n", err.Error())
					}
				}
			case sig := <-sigs:
				log.Printf("received signal: %s\n", sig.String())
				done <- true
			case _ = <-idle:
				log.Println("received idle signal")
				log.Println("sending shutdown signal")
				cmd := exec.Command("shutdown", "+1")
				if err = cmd.Run(); err != nil {
					log.Printf("shutdown error: %s\n", err.Error())
				}

				done <- true
			}
		}
	}()

	//block
	<-done

	// stop playback
	if err := mopidyClient.stop(); err != nil {
		fmt.Printf("mopidy.quit: %s", err.Error())
	}

	os.Exit(0)
}
