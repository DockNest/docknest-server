package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DockNest/docknest-server/pkg/docker"
	"github.com/DockNest/docknest-server/pkg/handlers"
)

type alias string

func main() {
	if err := os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock"); err != nil {
		log.Fatal(err)
	}

	docker.NewDockerClient()

	port := ":8891"
	log.Println("Listening on port ", port)
	http.HandleFunc("/shipyard", handlers.ShipyardCommand)

	log.Fatal(http.ListenAndServe(port, nil))
}
