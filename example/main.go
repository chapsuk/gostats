package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/chapsuk/gostats"
)

func main() {
	h := flag.String("host", "192.168.99.100", "DOCKER_HOST")
	max := flag.Int("i", 1000, "iterations count")
	sleep := flag.Int("s", 100, "sleep milliseconds after iteration")
	stop := flag.Int("b", 500, "skip iteration number")
	cont := flag.Int("c", 800, "continue iteration number")
	flag.Parse()

	stat, err := gostats.NewStatsd(fmt.Sprintf("%s:8125", *h), "Gostats.Example.")
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("Start write metrics.")
	stat.SetErrorHandler(func(err error) {
		log.Print(err)
	})

	dch := make(chan bool)
	tchan := time.NewTicker(1 * time.Second).C
	go func() {
		for {
			select {
			case <-tchan:
				stat.WriteCounter("Counter", 100+random(0, 200))
			case <-dch:
				return
			}
		}
	}()

	for i := 0; i < *max; i++ {
		stat.WriteTiming("Timing", 30+random(0, 20))

		if i == *stop {
			stat.Stop()
		} else if i == *cont {
			stat.Continue()
		}

		time.Sleep(time.Duration(*sleep) * time.Millisecond)
	}

	dch <- true
	log.Print("The end!")
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
