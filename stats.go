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
	ch     chan *event
	dn     chan bool
	logger golog.StandartLogger
	active bool
	mu     sync.Mutex
}

type event struct {
	key  string
	time int // milliseconds
}

// NewStatsd create statsd statistic writer
func NewStatsd(address, p string) (*Statsd, error) {
	c := statsd.NewStatsdClient(address, p)
	err := c.CreateSocket()
	if err != nil {
		return nil, err
	}
	s := &Statsd{
		ch:     make(chan *event),
		dn:     make(chan bool),
		cli:    c,
		active: true,
	}
	go s.collect()
	return s, nil
}

// SetLogger set logger
func (s *Statsd) SetLogger(l golog.StandartLogger) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logger = l
}

// Write stats key, t is milliseconds
func (s *Statsd) Write(k string, t int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.ch <- &event{key: k, time: t}
	}
}

// Stop collect metrics
func (s *Statsd) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.dn <- true
		s.active = false
	}
}

// Continue collect metrics
func (s *Statsd) Continue() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active {
		go s.collect()
		s.active = true
	}
}

func (s *Statsd) log(e error) {
	if s.logger != nil {
		s.logger.Print(e)
	}
}

func (s *Statsd) collect() {
	for {
		select {
		case e := <-s.ch:
			err := s.cli.Timing(fmt.Sprintf(".%s.time", e.key), int64(e.time))
			if err != nil {
				s.log(err)
			}
			err = s.cli.Incr(fmt.Sprintf(".%s.count", e.key), 1)
			if err != nil {
				s.log(err)
			}
		case <-s.dn:
			return
		}
	}
}
