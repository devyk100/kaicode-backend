package worker

import (
	"fmt"
	"log"
)

func (w *Worker) createPythonFile(code string) string {
	pyFileName := "/tmp/python/program.py"
	createFileCmd := fmt.Sprintf(`cat << 'EOF' > %s
%s
EOF
`, pyFileName, code)
	_, err := w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {
		// Handle error
	}
	return pyFileName
}

func (w *Worker) execPython(testcaseInput string, filename string) chan string {
	c := make(chan string)
	go func() {
		runCmd := fmt.Sprintf("echo '%s' | python3 %s", testcaseInput, filename)
		runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
		fmt.Print(runOutput, "IS THE OUTPUT FROM PYTHON EXEC")
		if err != nil {
			// Handle error
			c <- err.Error()
		}
		c <- runOutput
	}()
	return c
}

func (w *Worker) cleanUpPython(filename string) {
	cleanupCmd := fmt.Sprintf("rm %s", filename)
	_, err := w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
