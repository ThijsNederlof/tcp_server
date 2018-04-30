[![Go Report Card](https://goreportcard.com/badge/github.com/thijsnederlof/tcp_server)](https://goreportcard.com/report/github.com/thijsnederlof/tcp_server)

# TCPServer
Fork of github.com/firstrow/tcp_server, but instead of receiving strings it uses byte arrays. 

### Install package

``` bash
go get -u github.com/thijsnederlof/tcp_server
```

### Usage:

NOTICE: This library does not use any delimiter to separate received messages. You will need to build your own
        framing protocol on top of this lib. 

``` go
package main

import "github.com/thijsnederlof/tcp_server"

func main() {
	server := tcp_server.New("localhost:9999")

	server.OnNewClient(func(c *tcp_server.Client) {
		// new client connected
		// lets send some message
		c.Send("Hello")
	})
	server.OnNewMessage(func(c *tcp_server.Client, message []byte) {
		// new message received
		// Frame messages here
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
	})

	server.Listen()
}
```
