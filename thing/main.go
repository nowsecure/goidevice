package main

import (
	"log"
	"os"

	"github.com/nowsecure/goidevice/idevice"
	"github.com/nowsecure/goidevice/lockdown"
)

func main() {
	device, err := idevice.New(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	lock, err := lockdown.NewClientWithHandshake(device, "thingy")
	if err != nil {
		log.Fatal(err)
	}
	client, err := lock.StartService(device, lockdown.CRASH_REPORT_MOVER_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	err = client.ReadPing()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("yay we did it")
	}
}
