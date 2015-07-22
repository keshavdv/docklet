package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"bytes"
	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/context"
	"github.com/nu7hatch/gouuid"
	"github.com/robfig/config"
	"github.com/unrolled/render"
)

var docker_client *docker.Client
var etcd_client *etcd.Client

const JOB_PENDING = 0
const JOB_SUCCESS = 1
const JOB_FAILURE = 2

var jobs map[string]int

func init() {
	var err error
	c, _ := config.ReadDefault("./config/docklet.conf")
	machines := []string{"http://localhost:4001"}
	etcd_client = etcd.NewClient(machines)
	jobs = make(map[string]int)

	docker_host, _ := c.String("docker", "host")
	docker_port, _ := c.String("docker", "port")
	endpoint := fmt.Sprintf("tcp://%s:%s", docker_host, docker_port)
	path := os.Getenv("DOCKER_CERT_PATH")
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	docker_client, err = docker.NewTLSClient(endpoint, cert, key, ca)
	if err != nil {
		log.Panic("Could not connect to docker!")
	}
}

func pullImage(image string) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	jobs[u.String()] = JOB_PENDING
	go func() {
		log.Println("TETS")
		var buf bytes.Buffer
		buf.Reset()
		err := docker_client.PullImage(docker.PullImageOptions{Repository: image, OutputStream: &buf}, docker.AuthConfiguration{})
		log.Println(buf.String())
		if err != nil {
			log.Println(err.Error())
			jobs[u.String()] = JOB_FAILURE
		} else {
			jobs[u.String()] = JOB_SUCCESS
		}
		log.Println("DONE")
	}()
	return u.String(), nil
}

func Home(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)
	r.HTML(w, http.StatusOK, "index", nil)
}

func Build(w http.ResponseWriter, req *http.Request) {
	// TODO: handle Dockerfile builds
}

func Pull(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)

	image := req.URL.Query().Get("image")
	if image == "" {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "invalid", "msg": "image must be specified"})
		return
	}

	var buf bytes.Buffer
	buf.Reset()
	docker_client.PullImage(docker.PullImageOptions{Repository: image, OutputStream: &buf}, docker.AuthConfiguration{})
	log.Println(buf.String())
}

func Create(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)

	image := req.URL.Query().Get("image")
	if image == "" {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "invalid", "msg": "image must be specified"})
		return
	}

	container, err := docker_client.CreateContainer(docker.CreateContainerOptions{Config: &docker.Config{OpenStdin: true, StdinOnce: true, AttachStdin: true, AttachStderr: true, AttachStdout: true, Image: image, Cmd: []string{"/bin/bash"}, Tty: true}})
	if err != nil {
		if err == docker.ErrNoSuchImage {
			// attempt to pul
			job, err := pullImage(image)
			if err != nil {
				r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed"})
			}
			r.JSON(w, http.StatusOK, map[string]string{"status": "pending", "job": job})
		} else {
			r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed", "msg": err.Error()})
		}
		return
	}
	r.JSON(w, http.StatusOK, map[string]string{"status": "created", "id": container.ID})
}

func Inspect(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)
	// inspect container to find what host ports are exposed
	id := req.URL.Query().Get("id")
	if id == "" {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "invalid", "msg": "job id must be specified"})
		return
	}
	container_info, err := docker_client.InspectContainer(id)
	if err != nil {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed", "msg": err.Error()})
		return
	}
	r.JSON(w, http.StatusOK, map[string]*docker.Container{"info": container_info})
}

func Status(w http.ResponseWriter, req *http.Request) {
	// Status of long running tasks like pull/build
	r := context.Get(req, "Render").(*render.Render)

	id := req.URL.Query().Get("job_id")
	if id == "" {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "invalid", "msg": "job id must be specified"})
		return
	}
	if val, ok := jobs[id]; ok {
		var status string

		if val == JOB_SUCCESS {
			status = "complete"
		} else if val == JOB_FAILURE {
			status = "failure"
		} else {
			status = "pending"
		}

		r.JSON(w, http.StatusOK, map[string]string{"status": status})
	} else {
		r.JSON(w, http.StatusNotFound, map[string]string{"status": "unknown"})

	}
}

func Start(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)

	id := req.URL.Query().Get("id")
	if id == "" {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "invalid", "msg": "id must be specified"})
		return
	}

	err := docker_client.StartContainer(id, &docker.HostConfig{PublishAllPorts: true})
	if err != nil {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed", "msg": err.Error()})
		return
	}

	// inspect container to find what host ports are exposed
	container_info, err := docker_client.InspectContainer(id)
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
