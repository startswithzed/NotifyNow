package main

import "context"

// A Subscriber subscribes to a Publisher
// and listens to the messages sent by it.
type Subscriber struct {
  id     uint
  out    chan string
  ctx    context.Context
  cancel context.CancelFunc
}

// NewSubscriber returns a new Subscriber 
// with the specified id and an out channel 
// that receives the messages sent to it by its publisher.
func NewSubscriber(id uint, ctx context.Context, cancel context.CancelFunc) *Subscriber {
  return &Subscriber{id: id, out: make(chan string), ctx: ctx, cancel: cancel}
}