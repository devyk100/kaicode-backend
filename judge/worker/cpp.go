package worker

import (
	"fmt"
	"log"
)

func (w *Worker) createCppFile(code string) string {
	cppFileName := "/tmp/cpp/program.cpp"
	createFileCmd := fmt.Sprintf(`cat << 'EOF' > %s
%s
EOF
`, cppFileName, code)
	_, err := w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {

	}
	return cppFileName
}

func (w *Worker) compileCpp(filename string) (string, error) {
	compileCmd := fmt.Sprintf("g++ %s -o /tmp/cpp/program", filename)
	compileOutput, err := w.dockerContainer.ExecInContainerStdCopy(compileCmd)
	if compileOutput != "" {
		err = fmt.Errorf("compile error")
	}
	return compileOutput, err
}

func (w *Worker) execCpp(testcaseInput string) chan string {
	c := make(chan string)
	go func() {
		runCmd := fmt.Sprintf("echo '%s' > /tmp/input.txt && /tmp/cpp/program < /tmp/input.txt", testcaseInput)

		runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
		if err != nil {
			c <- err.Error()
			return
		}
		c <- runOutput
	}()
	return c
}

func (w *Worker) cleanUpCpp(filename string) {
	cleanupCmd := fmt.Sprintf("rm %s /tmp/cpp/program", filename)
	_, err := w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
