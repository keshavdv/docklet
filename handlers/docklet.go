package handlers

import (
	"net/http"
	"os"
	"fmt"
	"log"

	"github.com/coreos/go-etcd/etcd"
	"github.com/unrolled/render"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/context"
)


var docker_client *docker.Client
var etcd_client *etcd.Client

func init() {
	var err error
	machines := []string{"http://localhost:4001"}
	etcd_client = etcd.NewClient(machines)

	endpoint := "tcp://192.168.99.100:2376"
	path := os.Getenv("DOCKER_CERT_PATH")
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	docker_client, err = docker.NewTLSClient(endpoint, cert, key, ca)
	if err != nil {
		log.Panic("Could not connect to docker!")
	}
}

func Launch(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)

	image := req.URL.Query().Get("image")
	if image == "" {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "invalid", "msg": "image must be specified"})
		return
	}

	// create container
	container, err := docker_client.CreateContainer(docker.CreateContainerOptions{Config: &docker.Config{Cmd: []string{"/bin/sh", "-c", "sleep 10000000"}, Image: image}})
	if err != nil {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed", "msg": err.Error()})
		return
	}

	// start container
	log.Println(fmt.Sprintf("starting container (%s) from image (%s)", container.ID, image))
	err = docker_client.StartContainer(container.ID, &docker.HostConfig{PublishAllPorts: true})
	if err != nil {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed", "msg": err.Error()})
		return
	}

	// inspect container to find what host ports are exposed
	container_info, err := docker_client.InspectContainer(container.ID)
	if err != nil {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed", "msg": err.Error()})
		return
	}
	log.Println(container_info.NetworkSettings.Ports)
	for port, binding := range container_info.NetworkSettings.Ports {
		log.Println(port)
		log.Println(binding[0].HostIP)
		log.Println(binding[0].HostPort)
	}

	// TODO: update etcd for confd with port mappings

	// return all the relevant bits
	r.JSON(w, http.StatusOK, map[string]string{"status": "launched"})
}

