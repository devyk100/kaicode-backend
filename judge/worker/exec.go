package worker

import (
	"fmt"
	"kc-backend/judge/types"
	"time"
)

func (w *Worker) ExecCode(key string, job types.Job) types.FinishedPayload {
	fmt.Printf("Running job for key %s\n", key)
	testcaseInput := job.Input

	// compile and filecreation setup for all the languages
	var filename string
	var compileOut string
	var err error

	switch job.Language {
	case "cpp":
		filename = w.createCppFile(job.Code)
		compileOut, err = w.compileCpp(filename)
		defer w.cleanUpCpp(filename)
	case "java":
		filename = w.createJavaFile(job.Code)
		compileOut, err = w.compileJava(filename)
		defer w.cleanUpJava(filename)
	case "python":
		filename = w.createPythonFile(job.Code)
		defer w.cleanUpPython(filename)
	case "javascript":
		filename = w.createJavascriptFile(job.Code)
		defer w.cleanUpJavascript(filename)
	case "go":
		filename = w.createGoFile(job.Code)
		defer w.cleanUpGo(filename)
	default:
	}
	if err != nil {
		return types.FinishedPayload{
			SubmissionId: job.SubmissionId,
			Output:       compileOut,
			TimeTaken:    0,
			SQSKey:       key,
		}
	}

	var outputChan chan string
	var outputString string

	start := time.Now()
	switch job.Language {
	case "cpp":
		outputChan = w.execCpp(testcaseInput)
	case "java":
		outputChan = w.execJava(testcaseInput)
	case "python":
		outputChan = w.execPython(testcaseInput, filename)
	case "javascript":
		outputChan = w.execJavascript(testcaseInput, filename)
	case "go":
		outputChan = w.execGo(testcaseInput)
	default:
		return types.FinishedPayload{
			SubmissionId: job.SubmissionId,
			Output:       "invalid language",
		}
	}

	select {
	case <-time.After(1 * time.Minute):
		w.dockerContainer.RestartContainer()
		return types.FinishedPayload{
			SubmissionId: job.SubmissionId,
			Output:       "Your code took too long to execute",
			TimeTaken:    0,
			SQSKey:       key,
		}
	case outputString = <-outputChan:
	}

	since := time.Since(start)

	fmt.Println("Actual output", outputString)
	outputString = removeNonPrintableChars(outputString)

	fmt.Print("It all executed without any issues!!!")
	return types.FinishedPayload{
		SubmissionId: job.SubmissionId,
		Output:       outputString,
		TimeTaken:    int32(since.Milliseconds()),
		SQSKey:       key,
	}
}
