package tcp_server

import (
	"bufio"
	"log"
	"net"
)

// Client holds info about connection
type Client struct {
	conn   net.Conn
	Server *server
}

// TCP server
type server struct {
	address                  string // Address to open connection: localhost:9999
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message []byte)
}

// Read client data from channel
func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)
	buf := make([]byte, 2024)
	for {
		nbyte, err := reader.Read(buf)
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}

		// Copy and pass only what is read from the conn instead of entire buffer
		message := make([]byte, nbyte)
		copy(message, buf[:nbyte])

		c.Server.onNewMessage(c, message)
	}
}

// Send sends a string message to client
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// SendBytes sends a byte array to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

// Conn returns the connection from the client struct
func (c *Client) Conn() net.Conn {
	return c.conn
}

// Close closes the connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// Called right after server starts listening new client
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) OnNewMessage(callback func(c *Client, message []byte)) {
	s.onNewMessage = callback
}

// Listen opens the tcp server to accept new connections
func (s *server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.listen()
		s.onNewClientCallback(client)
	}
}

// New creates a new tcp server instance on specified address
func New(address string) *server {
	log.Println("Creating server with address", address)
	server := &server{
		address: address,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message []byte) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}
