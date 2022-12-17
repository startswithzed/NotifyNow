package main

import (
	"context"
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
	ctx         context.Context
	cancel      context.CancelFunc // its okay to pass functions as values
}

// NewPublisher returns a new Publisher with zero subscribers.
// A Publishers needs to be started with the Publisher.Start()
// method before it can start publishing incoming messages.
func NewPublisher() *Publisher {
	ctx, cancel := context.WithCancel(context.Background())
	return &Publisher{subscribers: map[uint]*Subscriber{}, ctx: ctx, cancel: cancel}
}

// AddSubscriber creates a new Subscriber that starts
// listening to the messages sent by this Publisher.
func (p *Publisher) AddSubscriber() {
	dctx, cancel := context.WithCancel(p.ctx)
	p.Lock()
	nextId := p.sequence + 1
	subscriber := NewSubscriber(nextId, dctx, cancel)
	p.subscribers[nextId] = subscriber
	p.sequence = p.sequence + 1
	p.Unlock()

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				fmt.Printf("parent context cancelled\n")
				subscriber.cancel()
				return
			case <-dctx.Done():
				fmt.Printf("subscriber %d's context cancelled\n", subscriber.id)
				return
			default:
				// check for the cancellation of the context before receiving the message
				select {
				case <-dctx.Done():
					fmt.Printf("subscriber %d's context cancelled\n", subscriber.id)
					return
				case msg := <-subscriber.out:
					fmt.Printf("message received on channel %d: %s\n", subscriber.id, msg)
				}
			}
		}
	}()

}

// Publish sends the msg to all the subscribers subscribed to it.
func (p *Publisher) Publish(msg string) error {
	// publish only if Publisher has been started
	if p.in == nil {
		return errors.New("publisher not started yet")
	}
	// publish only if subscriber count > 0
	p.RLock()
	if len(p.subscribers) == 0 {
		return errors.New("no subscribers to receive the message")
	}
	p.RUnlock()

	// send message to in channel
	// establish lock to ensure no modifications take place
	// while the message is being sent
	p.Lock()
	p.in <- msg
	p.Unlock()

	return nil
}

// Start initializes the the publishers in channel
// and makes it ready to start publishing incoming messages.
func (p *Publisher) Start() {
	in := make(chan string)
	p.in = in

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				// using break here only breaks out of select statement
				// instead of breaking out of the loop
				fmt.Printf("done called on publisher\n")
				return
			case msg := <-p.in:
				p.RLock()
				for _, subscriber := range p.subscribers {
					subscriber.out <- msg
				}
				p.RUnlock()
			}
		}
	}()
}

// Stop prevents the Publisher from listening to any
// incoming messages by closing the in channel.
func (p *Publisher) Stop() {
	p.cancel()
}

// RemoveSubscriber removes the subscriber specified by the given id.
// It returns an error "could not find subscriber"
// if Subscriber with the given id is not subscribed to the Publisher.
func (p *Publisher) RemoveSubscriber(id uint) error {
	p.Lock()
	defer p.Unlock()
	subscriber, ok := p.subscribers[id]

	if !ok {
		return errors.New("could not find subscriber")
	}

	subscriber.cancel()
	delete(p.subscribers, id)
	close(subscriber.out)

	return nil
}
