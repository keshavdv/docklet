package handlers

import (
	"bufio"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/context"
	"github.com/gorilla/websocket"
	"github.com/unrolled/render"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func Attach(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)
	r.HTML(w, http.StatusOK, "terminal", nil)
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type message struct {
	data []byte
	conn *connection
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump(writer io.Writer) {
	defer func() {
		Hub.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		writer.Write(message)
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serverWs handles websocket requests from the peer.
func CreateTerminalServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	Hub.register <- c

	go c.writePump()
	attachToContainer("e41e5d77750a", c)

}

func attachToContainer(containerId string, c *connection) {

	containerOutR, containerOut := io.Pipe()
	containerIn, containerInW := io.Pipe()

	go docker_client.AttachToContainer(docker.AttachToContainerOptions{
		Logs:         true,
		Stream:       true,
		Stdin:        true,
		Stdout:       true,
		Stderr:       true,
		Container:    containerId,
		InputStream:  containerIn,
		OutputStream: containerOut,
		ErrorStream:  containerOut,
		RawTerminal:  true,
	})

	go func(reader io.Reader, c *connection) {

		defer func() {
			containerOutR.Close()
			containerOut.Close()
			containerIn.Close()
			containerInW.Close()
		}()

		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanBytes)
		for scanner.Scan() {

			Hub.emit <- message{data: []byte(scanner.Text()), conn: c}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error with the scanner in attached container", err)
		}
	}(containerOutR, c)

	c.readPump(containerInW)
}
