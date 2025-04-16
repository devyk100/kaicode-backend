package worker

import (
	"fmt"
	"log"
)

func (w *Worker) createGoFile(code string) string {
	goFileName := "/tmp/cpp/program.go"
	createFileCmd := fmt.Sprintf(`cat << 'EOF' > %s
%s
EOF
`, goFileName, code)
	_, err := w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {

	}
	return goFileName
}

// compileGo deprecated, no need to compile the file before running it directly in golang, just use `go run file_name.go` directly
func (w *Worker) compileGo(filename string) (string, error) {
	compileCmd := fmt.Sprintf("GOCACHE=/tmp/go/go-cache go build -o /tmp/go/program %s", filename)
	compileOutput, err := w.dockerContainer.ExecInContainer(compileCmd)
	if compileOutput != "" {
		err = fmt.Errorf("compile error")
	}
	return compileOutput, err
}

func (w *Worker) execGo(testcaseInput string) chan string {
	c := make(chan string)
	go func() {
		runCmd := fmt.Sprintf("echo '%s' | GOCACHE=/tmp/go/go-cache go run /tmp/go/program.go", testcaseInput) // it was just "go run /tmp/go/program"
		runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
		if err != nil {
			c <- err.Error()
		}
		c <- runOutput
	}()
	return c
}

func (w *Worker) cleanUpGo(filename string) {
	cleanupCmd := fmt.Sprintf("rm %s /tmp/go/program", filename)
	_, err := w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
