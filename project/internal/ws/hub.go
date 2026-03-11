package ws

type BroadcastMsg struct {
	UserIDs []string
	Data    []byte
}
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMsg
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMsg),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if old, ok := h.clients[client.UserID]; ok {
				close(old.send)
			}
			h.clients[client.UserID] = client
		case client := <-h.unregister:
			if cl, ok := h.clients[client.UserID]; ok && client == cl {
				delete(h.clients, client.UserID)
				close(client.send)
			}
		case msg := <-h.broadcast:
			for _, UserID := range msg.UserIDs {
				if client, ok := h.clients[UserID]; ok {
					select {
					case client.send <- msg.Data:
					default:
						delete(h.clients, UserID)
					}
				}
			}
		}
	}
}

func (h *Hub) Register(c *Client) {
	h.register <- c
}
func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}
func (h *Hub) Broadcast(msg *BroadcastMsg) {
	h.broadcast <- msg
}
