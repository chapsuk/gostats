package gostats

import (
	"fmt"
	"sync"

	"github.com/quipo/statsd"
)

// ErrorHandler is error handler function
type ErrorHandler func(err error)

// Statsd metrics writer
type Statsd struct {
	cli    *statsd.StatsdClient
	tch    chan *timingEvent
	tdn    chan bool
	cch    chan *countEvent
	cdn    chan bool
	active bool
	mu     sync.Mutex
	eh     ErrorHandler
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

// SetErrorHandler set error handler function
func (s *Statsd) SetErrorHandler(f ErrorHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eh = f
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

func (s *Statsd) log(e error) {
	if s.eh != nil {
		s.eh(e)
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
			s.eh(fmt.Errorf("stop collect timings"))
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
			s.eh(fmt.Errorf("stop collect counters"))
			return
		}
	}
}
