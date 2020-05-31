package main

import "fmt"
import "log"
import "strconv"
import "encoding/hex"
import "github.com/jacobsa/go-serial/serial"

func main() {
	options := serial.OpenOptions{
      PortName: "/dev/ttyS0",
      BaudRate: 9600,
      DataBits: 8,
      StopBits: 1,
      MinimumReadSize: 14,
    }

    // Open the port.
    port, err := serial.Open(options)
    if err != nil {
      log.Fatalf("serial.Open: %v", err)
    }

    // Make sure to close it later.
    defer port.Close()

        buf := make([]byte, 16)

        for {
                r, err := port.Read(buf)
                if err != nil {
                        log.Fatal(err)
                }
                s := hex.Dump(buf[:r])
		s2 := string(buf[4:11])
                fmt.Print(s)
		fmt.Println("")
	n, err := strconv.ParseUint(s2, 16, 32)
if err != nil {
   continue 
}
n2 := uint32(n)
		fmt.Printf("%d",n2)	
		fmt.Println("")
        }
}
