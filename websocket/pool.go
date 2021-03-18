package websocket

import "fmt"

// Pool struct will contain all of the channels we need for concurrent communication,
// as well as a map of clients.
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

// `Start()` method will constantly listen for anything passed to any of our Pool’s
// channels and then, if anything is received into one of these channels,
// it’ll act accordingly
func (pool *Pool) Start() {
	for {
		select {
		// `Register` channel will send out New User Joined... to all of the clients
		// within this pool when a new client connects.
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
			break
		// `Unregister` will unregister a user and notify the pool when a client disconnects.
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
			break
		// `Broadcast` a channel which, when it is passed a message,
		// will loop through all clients in the pool and send the message
		// through the socket connection.
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
