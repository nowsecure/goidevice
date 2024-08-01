package main

import (
	"fmt"
	"log"

	"github.com/nowsecure/goidevice/afc"
	"github.com/nowsecure/goidevice/idevice"
	"github.com/nowsecure/goidevice/lockdown"
)

func main() {
	device, err := idevice.New("bd133240a37062e545bbbbf664f0011c9f45895d")
	if err != nil {
		log.Fatal(err)
	}
	lock, err := lockdown.NewClientWithHandshake(device, "thingy")
	if err != nil {
		log.Fatal(err)
	}
	client, err := lock.StartServiceClient(device, lockdown.CRASH_REPORT_MOVER_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Free()
	err = client.ReadPing()
	if err != nil {
		log.Fatal(err)
	}

	service, err := lock.StartService(device, lockdown.CRASH_REPORT_COPY_MOBILE_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Free()

	afc, err := afc.NewClient(device, service)
	if err != nil {
		log.Fatal(err)
	}
	k, _ := afc.WalkDirectory(".")
	for _, v := range k {
		fmt.Println(v)
	}

	log.Println("yay we did it")
}
