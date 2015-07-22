package handlers

import (
	"github.com/googollee/go-socket.io"
	"log"
	"io"
	"github.com/fsouza/go-dockerclient"
	"fmt"
	"bufio"
	"os"
	"net/http"
	"github.com/gorilla/context"
	"github.com/unrolled/render"
)

func Attach(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)
	r.HTML(w, http.StatusOK, "terminal", nil)
}

func CreateTerminalServer() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		containerId := "93a7eaccb953"
		attachToContainer(so, containerId)
		log.Println("on connection")

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return server
}

func attachToContainer(so socketio.Socket, containerId string) {

	containerOutR, containerOut := io.Pipe()
	containerIn, containerInW := io.Pipe()

	go docker_client.AttachToContainer(docker.AttachToContainerOptions{
		Logs: true,
		Stream: true,
		Stdin: true,
		Stdout: true,
		Stderr: true,
		Container: containerId,
		InputStream: containerIn,
		OutputStream: containerOut,
		ErrorStream: containerOut,
		RawTerminal: true,
	})

	go func(reader io.Reader, so socketio.Socket) {
		fmt.Println("in go routine")

		defer containerOutR.Close()
		defer containerOut.Close()
		defer containerIn.Close()
		defer containerInW.Close()

		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanBytes)
		for scanner.Scan() {

			if so != nil {
				err := so.Emit("output", fmt.Sprintf("%s", scanner.Text()))
				if err != nil {
					return
				}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error with the scanner in attached container", err)
		}
	}(containerOutR, so)

	so.On("input", func(msg string) {
		containerInW.Write([]byte(msg))
	})

}