package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	room *Room
	conn *websocket.Conn
	send chan []byte
}

func NewClient(room *Room, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	client := &Client{
		room: room,
		conn: conn,
		send: make(chan []byte, 256),
	}
	go client.write()
	go client.read()
}

func (this *Client) read() {

}

func (this *Client) write() {

}
