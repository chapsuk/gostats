package main

import (
	"flag"
	"fmt"
	"github.com/chapsuk/gostats"
	"log"
	"time"
)

func main() {
	h := flag.String("h", "192.168.99.100", "DOCKER_HOST")
	flag.Parse()
	stat, err := gostats.NewStatsd(fmt.Sprintf("%s:8125", *h), "gostat.test")
	if err != nil {
		log.Print(err)
		return
	}
	t := time.Now()
	for i := 0; i < 1000; i++ {
		ct := int(time.Now().Sub(t)) / int(time.Millisecond)
		log.Printf("Write %d metric, time: %d ms", i, ct)
		stat.Write("foo.bar", ct)
		if i == 500 {
			stat.Stop()
		} else if i == 800 {
			stat.Continue()
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Print("The end!")
}
