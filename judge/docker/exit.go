package docker

import (
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
)

func (d *Docker) Exit() {
	fmt.Println("Exiting the docker container")
	err := d.cli.ContainerKill(d.ctx, d.containerId, "SIGKILL")
	if err != nil {
		log.Fatal("Error inspecting container: ", err)
	}
	err = d.cli.ContainerRemove(d.ctx, d.containerId, container.RemoveOptions{Force: true})
	if err != nil {
		log.Printf("Error removing container %s: %v", d.containerId, err)
	} else {
		fmt.Printf("Container %s removed successfully.\n", d.containerId)
	}

}
