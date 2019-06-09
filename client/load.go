package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/howardjohn/pilot-load/adsc"
)

func makeADSC(addr string, client int, prefix int) error {
	ip := fmt.Sprintf("%d.0.%d.%d", prefix, client/256, client%256)
	fmt.Println("Connecting:", ip)
	con, err := adsc.Dial(addr, "", &adsc.Config{
		IP: ip,
	})
	if err != nil {
		return err
	}
	fmt.Println("Connected:", ip)
	con.Watch()
	fmt.Println("Got Initial Update:", ip)
	for {
		u := <-con.Updates
		if u == "close" {
			fmt.Println("Closing:", ip)
			return nil
		}
	}
}

func RunLoad(pilotAddress string, clients int, prefix int) error {
	wg := sync.WaitGroup{}
	for cur := 0; cur < clients; cur++ {
		wg.Add(1)
		go func() {
			makeADSC(pilotAddress, cur, prefix)
			wg.Done()
		}()
		time.Sleep(time.Millisecond * 100)
	}
	wg.Wait()
	return nil
}