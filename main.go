package main

import (
	"fmt"
	"net/http"
)

func main() {
  var err error
  
	publisher := NewPublisher()

  publisher.AddSubscriber()
  publisher.AddSubscriber()
  publisher.AddSubscriber()

  publisher.Start()
  
  err = publisher.Publish("Hello World")
  if err != nil {
    fmt.Printf("could not publish: %v\n", err)
  }

  err = publisher.RemoveSubscriber(1)
  if err != nil {
    fmt.Printf("could not remove subscriber: %v\n", err)
  }
  
  err = publisher.Publish("Only channels 2 and 3 should print this")
  if err != nil {
    fmt.Printf("could not publish: %v\n", err)
  }

  // this is merely to keep the server running
  http.ListenAndServe(":8080", nil)
}
