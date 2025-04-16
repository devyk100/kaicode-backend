package docker

func (d *Docker) RestartContainer() {
	d.Exit()
	err := d.StartContainer(d.ctx)
	if err != nil {
		// configure the appropriate behaviour
	}
}
