package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

var cli *client.Client

type DockerCmdOptions struct {
	ContainerID string
	Limit       uint
}

func NewDockerClient() {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal("could not create docker client ", err.Error())
	}

	cli = c
}

var Commands = map[string]func(opts DockerCmdOptions) ([]byte, error){
	"ps":   showRunningContainers,
	"logs": showLogs,
}

func showRunningContainers(o DockerCmdOptions) ([]byte, error) {
	cs, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return json.Marshal(cs)
}

func showLogs(o DockerCmdOptions) ([]byte, error) {
	if o.ContainerID == "" {
		log.Fatal("container id is empty")
	}

	tail := "32"
	if o.Limit > 0 && o.Limit <= 4096 {
		tail = strconv.Itoa(int(o.Limit))
	}

	log.Printf("limit: %s", tail)

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	}

	logs, err := cli.ContainerLogs(context.Background(), o.ContainerID, options)
	if err != nil {
		log.Fatalf("failed to fetch container logs: %v", err)
	}

	defer logs.Close()

	var outBuf, errBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuf, &errBuf, logs)
	if err != nil {
		log.Fatalf("failed to demultiplex container logs: %v", err)
	}

	allLogs := append(outBuf.Bytes(), errBuf.Bytes()...)

	jsonLogs, err := json.Marshal(strings.Split(string(allLogs), "\n"))
	if err != nil {
		log.Fatalf("failed to marhshal logs into json %v", err)
	}

	return jsonLogs, nil
}
