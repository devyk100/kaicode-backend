package worker

func (w *Worker) Exit() {
	// explicitly exit this function
	w.IsRunning = false
	w.dockerContainer.Exit()
	// cancel the context
}
