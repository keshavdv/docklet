package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/context"
	"github.com/nu7hatch/gouuid"
	"github.com/robfig/config"
	"github.com/unrolled/render"
)

var docker_client *docker.Client
var etcd_client *etcd.Client

const JOB_PENDING = "pending"
const JOB_SUCCESS = "success"
const JOB_FAILURE = "failed"

var jobs map[string]*Job

type Job struct {
	Status string            `json:"status"`
	Data   map[string]string `json:"data,string"`
}

func init() {
	var err error
	c, _ := config.ReadDefault("./config/docklet.conf")
	machines := []string{"http://localhost:4001"}
	etcd_client = etcd.NewClient(machines)
	jobs = make(map[string]*Job)

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

func createContainer(image string, cmd []string) (*docker.Container, error) {
	container, err := docker_client.CreateContainer(
		docker.CreateContainerOptions{
			Config: &docker.Config{
				OpenStdin:    true,
				StdinOnce:    true,
				AttachStdin:  true,
				AttachStderr: true,
				AttachStdout: true,
				Image:        image,
				Cmd:          cmd,
				Tty:          true,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	log.Println("created")

	return container, nil
}

func pullImage(image string) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	go func() {
		jobs[u.String()] = &Job{JOB_PENDING, nil}
		container, err := createContainer(image, []string{"/bin/bash"})
		if err != nil {
			if err == docker.ErrNoSuchImage {
				log.Println("I need to pull image")
				err = docker_client.PullImage(docker.PullImageOptions{Repository: image}, docker.AuthConfiguration{})
				if err != nil {
					log.Println(err.Error())
					jobs[u.String()] = &Job{JOB_FAILURE, map[string]string{"message": err.Error()}}
				} else {
					container, err = createContainer(image, []string{"/bin/bash"})
					if err != nil {
						jobs[u.String()] = &Job{JOB_FAILURE, map[string]string{"message": err.Error()}}
					} else {
						jobs[u.String()] = &Job{JOB_SUCCESS, map[string]string{"id": container.ID}}
					}
				}
			} else {
				log.Println("failure")
				jobs[u.String()] = &Job{JOB_FAILURE, map[string]string{"message": err.Error()}}
			}
		} else {
			log.Println("hi")
			jobs[u.String()] = &Job{JOB_SUCCESS, map[string]string{"id": container.ID}}
		}
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

func Create(w http.ResponseWriter, req *http.Request) {
	r := context.Get(req, "Render").(*render.Render)

	image := req.URL.Query().Get("image")
	if image == "" {
		r.JSON(w, http.StatusBadRequest, map[string]string{"status": "invalid", "msg": "image must be specified"})
		return
	}

	job, err := pullImage(image)
	if err != nil {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed"})
	}
	r.JSON(w, http.StatusOK, map[string]string{"status": "pending", "job": job})

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
	if job, ok := jobs[id]; ok {
		r.JSON(w, http.StatusOK, job)
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
