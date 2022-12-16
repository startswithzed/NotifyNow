package main

import "net/http"

func main() {
	publisher := NewPublisher()

  publisher.AddSubscriber()
  publisher.AddSubscriber()
  publisher.AddSubscriber()

  publisher.Start()

  publisher.Publish("Hello World")

  publisher.RemoveSubscriber(1)
  publisher.Publish("Only channels 2 and 3 should print this")

  // publisher.Stop()

  // this is merely to keep the server running
  http.ListenAndServe(":8080", nil)
}
