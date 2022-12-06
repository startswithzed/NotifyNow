package main

import "net/http"

func main() {
	publisher := NewPublisher()

  publisher.AddSubscriber()
  publisher.AddSubscriber()
  publisher.AddSubscriber()

  publisher.Start()

  publisher.Publish("Hello World")

  // this is merely to keep the server running
  http.ListenAndServe(":8080", nil)
}
