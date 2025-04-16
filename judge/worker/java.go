package worker

import (
	"fmt"
	"log"
)

func (w *Worker) createJavaFile(code string) string {
	javaFileName := "/tmp/java/Main.java"
	createFileCmd := fmt.Sprintf(`cat << 'EOF' > %s
%s
EOF
`, javaFileName, code)

	_, err := w.dockerContainer.ExecInContainer(createFileCmd)
	if err != nil {
		fmt.Println("Create Java file error:", err.Error())
	}
	return javaFileName
}

func (w *Worker) compileJava(filename string) (string, error) {
	compileCmd := fmt.Sprintf("javac %s", filename)
	compileOutput, err := w.dockerContainer.ExecInContainer(compileCmd)
	fmt.Println(compileOutput)
	if compileOutput != "" {
		err = fmt.Errorf("compile error")
	}
	return compileOutput, err
}

func (w *Worker) execJava(testcaseInput string) chan string {
	c := make(chan string)
	go func() {
		runCmd := fmt.Sprintf("echo '%s' | java -cp /tmp/java Main", testcaseInput)
		runOutput, err := w.dockerContainer.ExecInContainer(runCmd)
		if err != nil {
			c <- err.Error()
			return
		}
		c <- runOutput
	}()
	return c
}

func (w *Worker) cleanUpJava(filename string) {
	cleanupCmd := fmt.Sprintf("rm %s /tmp/java/Program.class", filename)
	_, err := w.dockerContainer.ExecInContainer(cleanupCmd)
	if err != nil {
		log.Printf("Warning: failed to clean up files: %v", err)
	}
}
