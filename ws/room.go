package ws

type Room struct {
	clients    map[Client]bool
	broadcast  chan []byte
	register   chan Client
	unregister chan Client
}

func NewRoom() *Room {
	return &Room{
		clients:    make(map[Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan Client),
		unregister: make(chan Client),
	}
}

func (this *Room) Run() {
	for {
		select {
		case client := <-this.register:
			this.clients[client] = true
		case client := <-this.unregister:
			if _, ok := this.clients[client]; ok {
				close(client.send)
				delete(this.clients, client)
			}
		case message := <-this.broadcast:
			for client := range this.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(this.clients, client)
				}
			}
		}
	}
}
