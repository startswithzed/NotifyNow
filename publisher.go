package main

// todo add error checks

import (
	"errors"
	"fmt"
	"sync"
)

// Publisher sends messages received in its in channel
// to all the Subscribers listening to it.
type Publisher struct {
  sequence    uint
	in          chan string
	subscribers map[uint]*Subscriber
	sync.RWMutex
}

// NewPublisher returns a new Publisher with zero subscribers.
// A Publishers needs to be started with the Publisher.Start() 
// method before it can start publishing incoming messages.
func NewPublisher() *Publisher {
	return &Publisher{subscribers: map[uint]*Subscriber{}}
}

// AddSubscriber creates a new Subscriber that starts 
// listening to the messages sent by this Publisher.
func (p *Publisher) AddSubscriber() {
	p.Lock()
	nextId := p.sequence + 1
	subscriber := NewSubscriber(nextId)
	p.subscribers[nextId] = subscriber
  p.sequence = p.sequence + 1
	p.Unlock()
	
	go func() {
		for msg := range subscriber.out {
			fmt.Printf("message received on channel %d: %s\n", subscriber.id, msg)
		}
	}()
}

// Publish sends the msg to all the subscribers subscribed to it.
func (p *Publisher) Publish(msg string) {
	p.in <- msg
}

// Start initializes the the publishers in channel 
// and makes it ready to start publishing incoming messages.
func (p *Publisher) Start() {
	in := make(chan string)
	p.in = in

	go func() {
		for msg := range p.in {
			msgCopy := msg
			p.RLock()
			for _, subscriber := range p.subscribers {
				subscriber.out <- msgCopy
			}
			p.RUnlock()
		}
	}()
}

// Stop prevents the Publisher from listening to any
// incoming messages by closing the in channel.
func (p *Publisher) Stop() {
	close(p.in)
}

// RemoveSubscriber removes the subscriber specified by the given id.
// It returns an error "could not find subscriber"
// if Subscriber with the given id is not subscribed to the Publisher.
func(p *Publisher) RemoveSubscriber(id uint) error {
  subscriber, ok := p.subscribers[id]
  
  if !ok {
    return errors.New("could not find subscriber")
  }
  
  p.Lock()
  delete(p.subscribers, id)
  p.Unlock()
  
  close(subscriber.out)

  return nil
}

