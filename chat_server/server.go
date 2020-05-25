package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn net.Conn
}

// Send message sends a message to a specific client
func (c *Client) SendMsg(msg string) error {
	_, err := c.conn.Write([]byte(msg))
	return err
}

// Server is a TCP server which receives the messages sent by
// a particular user and relays it to all the other users
type server struct {
	address string

	// Handle concurrent access to the client
	// list from multiple goroutines
	mux sync.Mutex

	// All the clients to which
	// resend the messages
	clients []*Client
}

// New creates a new server with the specified address
func New(address string) *server {
	return &server{
		address: address,
		clients: []*Client{},
	}
}

func (s *server) handleClient(c *Client) {
	fmt.Printf("Handling client %s\n", c.conn.RemoteAddr().String())

	// Receive a message and send it to the others
	reader := bufio.NewReader(c.conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Got message! Resending to the others\n")

		// Resend the message to the other users
		s.mux.Lock()

		// Send the message to all the clients
		for _, client := range s.clients {
			if client != c {
				err = client.SendMsg(msg)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		s.mux.Unlock()
	}

}

// Listen listens for connections and handles each of them
// in a different goroutine
func (s *server) Listen() {

	// Create the listener
	fmt.Printf("Server listening in port %s\n", s.address)
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	// Serve the endpoint forever and
	// handle the client calls in a new routine
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		// Create the client object
		c := &Client{conn}

		// Add the client to our slice
		s.mux.Lock()
		s.clients = append(s.clients, c)
		s.mux.Unlock()

		// Handle the connection in other thread
		go s.handleClient(c)

	}

}

func main()  {

	// Create the server
	s := New("localhost:9999")

	// Channel to wait
	done := make(chan bool)

	// Make the server listen in another thread
	go s.Listen()

	// dial the server
	for i:= 0; i<3; i++ {

		// Start the different clients
		i := i
		go func(number int) {
			conn, err := net.Dial("tcp", "localhost:9999")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Client %d connected to the server\n", i)

			// Start a routine for printing tge messages
			go func(c net.Conn, number int) {
				for{
					msg , err := bufio.NewReader(c).ReadString('\n')
					if err != nil{
						panic(err)
					}
					fmt.Printf("Client %d got message %s\n", i, msg)

				}
			}(conn, i)

			time.Sleep(time.Second)

			fmt.Printf("Client %d writing to the server\n", i)
			conn.Write([]byte(fmt.Sprintf("Hello from client %d\n", i)))
		}(i)


	}

fmt.Println("Waiting...")
<- done


}
