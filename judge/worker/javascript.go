package worker

import (
	"fmt"
	"log"
)

func (w *Worker) createJavascriptFile(code string) string {
	jsFileName := "/tmp/js/program.js"
	createFileCmd := fmt.Sprintf(`cat << 'EOF' > %s
%s
EOF
`, jsFileName, code)
	_, err := w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {
		// Handle error
	}
	return jsFileName
}

func (w *Worker) execJavascript(testcaseInput string, filename string) chan string {
	c := make(chan string)
	go func() {
		runCmd := fmt.Sprintf("echo '%s' | node %s", testcaseInput, filename)
		runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
		fmt.Print(runOutput)
		if err != nil {
			// Handle error
			c <- err.Error()
		}
		c <- runOutput
	}()
	return c
}

func (w *Worker) cleanUpJavascript(filename string) {
	cleanupCmd := fmt.Sprintf("rm %s", filename)
	_, err := w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
