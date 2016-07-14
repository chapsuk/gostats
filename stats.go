package gostats

import (
	"fmt"
	"sync"

	"github.com/chapsuk/golog"
	"github.com/quipo/statsd"
)

// Statsd metrics writer
type Statsd struct {
	cli    *statsd.StatsdClient
	tch    chan *timingEvent
	tdn    chan bool
	cch    chan *countEvent
	cdn    chan bool
	logger golog.StandartLogger
	active bool
	mu     sync.Mutex
}

type timingEvent struct {
	key  string
	time int // milliseconds
}

type countEvent struct {
	key   string
	count int
}

// NewStatsd create statsd statistic writer
func NewStatsd(address, p string) (*Statsd, error) {
	c := statsd.NewStatsdClient(address, p)
	err := c.CreateSocket()
	if err != nil {
		return nil, err
	}
	s := &Statsd{
		tch:    make(chan *timingEvent),
		tdn:    make(chan bool),
		cch:    make(chan *countEvent),
		cdn:    make(chan bool),
		cli:    c,
		active: true,
	}
	go s.collectTimings()
	go s.collectCounters()
	return s, nil
}

// Close channels
func (s *Statsd) Close() {
	s.tdn <- true
	s.cdn <- true
	close(s.tch)
	close(s.tdn)
	close(s.cch)
	close(s.cdn)
	s.active = false
}

// SetLogger set logger
func (s *Statsd) SetLogger(l golog.StandartLogger) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logger = l
}

// WriteTiming write timing for key
func (s *Statsd) WriteTiming(k string, t int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.tch <- &timingEvent{key: k, time: t}
	}
}

// WriteCounter write counter metric with key
func (s *Statsd) WriteCounter(k string, c int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.cch <- &countEvent{key: k, count: c}
	}
}

// Stop collect metrics
func (s *Statsd) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.cdn <- true
		s.tdn <- true
		s.active = false
	}
}

// Continue collect metrics
func (s *Statsd) Continue() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active {
		go s.collectTimings()
		go s.collectCounters()
		s.active = true
	}
}

func (s *Statsd) log(m interface{}) {
	if s.logger != nil {
		s.logger.Print(m)
	}
}

func (s *Statsd) collectTimings() {
	for {
		select {
		case e := <-s.tch:
			err := s.cli.Timing(fmt.Sprintf(".%s.time", e.key), int64(e.time))
			if err != nil {
				s.log(err)
			}
		case <-s.tdn:
			s.log("stop collect timings")
			return
		}
	}
}

func (s *Statsd) collectCounters() {
	for {
		select {
		case e := <-s.cch:
			err := s.cli.Incr(fmt.Sprintf(".%s.count", e.key), int64(e.count))
			if err != nil {
				s.log(err)
			}
		case <-s.cdn:
			s.log("stop collect counters")
			return
		}
	}
}
