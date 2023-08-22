package handlers

import (
	"encoding/json"
	"github.com/DockNest/docknest-server/pkg/docker"
	"log"
	"net/http"
)

type RequestedCommand struct {
	Command     string // only works if command starts with a capital 'C'
	DockerCmd   string
	ContainerId string
	Limit       uint
	Test        bool
}

func ShipyardCommand(w http.ResponseWriter, r *http.Request) {
	var reqCmd RequestedCommand

	err := json.NewDecoder(r.Body).Decode(&reqCmd)
	if err != nil {
		http.Error(w, "Failed to parse req body", http.StatusBadRequest)
		return
	}

	log.Printf("Requested: %+v", reqCmd)

	if reqCmd.Test {
		w.Write([]byte("Honk honk!"))
		return
	}

	if cmd, ok := docker.Commands[reqCmd.DockerCmd]; ok {
		result, err := cmd(docker.DockerCmdOptions{
			ContainerID: reqCmd.ContainerId,
			Limit:       reqCmd.Limit,
		})

		if err != nil {
			w.Write([]byte("Whoops something went wrong executing the cmd..."))
		} else {
			log.Printf("Writing result: %s", string(result))
			w.Header().Add("Content-Type", "application/json")
			w.Write(result)
		}
	} else {
		log.Println("Unknown command: ", reqCmd.Command)
		w.Write([]byte("Command not supported"))
	}
}
