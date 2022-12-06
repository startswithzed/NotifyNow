package main

import (
	"fmt"
	"sync"
)

type Publisher struct {
  in           chan string
  subscribers  []chan <- string
  sync.RWMutex
}

func NewPublisher() *Publisher {
  return &Publisher{subscribers: []chan <- string{}}
}

func (p *Publisher) AddSubscriber() {
  out := make(chan string)
  p.Lock()
  p.subscribers = append(p.subscribers, out)
  p.Unlock()
  // start listening to this channel
  // todo: implement closing of this channel
  go func() {
    for msg:= range out {
      fmt.Println("message received: %s", msg)
    }
  }()
}

func (p *Publisher) Publish(msg string) {
  p.in <- msg
}

func (p *Publisher) Start() {
  in := make(chan string)
  p.in = in

  // todo: implement exiting of this goroutine    
  go func() {
    for msg := range p.in {
      msgCopy := msg
      p.RLock()
      for _, out := range p.subscribers {
        out <- msgCopy
      }
      p.RUnlock()
    }
  }()
}

