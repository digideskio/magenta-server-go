package main

import (
	"bufio"
	"net"
)

type Client struct {
	NickName    string
	RealName    string
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	serverInput chan Message
	output      chan string
}

func NewClient(name string, conn net.Conn, serverInput chan Message) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	client := &Client{
		NickName:    name,
		conn:        conn,
		reader:      reader,
		writer:      writer,
		serverInput: serverInput,
		output:      make(chan string),
	}

	// Start the loops that constantly monitor for input and output
	// from the client
	client.run()

	return client
}

// run starts two loops, concurrently, that monitor the Client's input and
// output channels
func (c *Client) run() {
	go c.processInput()
	go c.SendOutput()
}

// GetInput receives data from the Client and sends it to the server
func (c *Client) processInput() {
	for {
		// The buffer will stop reading when it detects a new line chatacter
		input, _ := c.reader.ReadString('\n')

		// If the client sends an empty message just ignore it
		if input := trimMessage(input); !isEmpty(input) {
			msg := NewMessage(c, input)
			c.serverInput <- *msg
		}
	}
}

// Close closes the clients connection and frees resources
func (c *Client) Close(msg string) {
	c.output <- msg
	c.conn.Close()
	close(c.output)
}

// SendOutput monitors the output channel.  When data is received from the server
// it is sent to the Client
func (c *Client) SendOutput() {
	for data := range c.output {
		c.writer.WriteString(data)
		c.writer.Flush()
	}
}
