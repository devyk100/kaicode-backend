package docker

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

func (d *Docker) ExecInContainer(command string) (string, error) {
	execConfig := container.ExecOptions{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execIDResp, err := d.cli.ContainerExecCreate(d.ctx, d.containerId, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %w", err)
	}

	resp, err := d.cli.ContainerExecAttach(d.ctx, execIDResp.ID, container.ExecAttachOptions{Tty: false})
	// WARNING: Do not use the simple io.ReadAll to read this buffer, you'll get fucked, as for large inputs the stdout has some random garbage characters that affect the comparison. To avoid that, you have to seperate out the stdout stderr as docker uses its own protocol, and that garbage is the result of that, so use docker's stdCopy method to extract stdout and stderr seperately. I found this out after 5-6 hours of debugging and finally coming across a 4 year old stackoverflow question, I am so stupid.
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer resp.Close()

	var stdoutBuf, stderrBuf bytes.Buffer

	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %w", err)
	}

	inspectResp, err := d.cli.ContainerExecInspect(d.ctx, execIDResp.ID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect exec instance: %w", err)
	}
	if inspectResp.ExitCode != 0 {
		return stdoutBuf.String() + "\n\n\nOperation failed with exit code: " + strconv.Itoa(inspectResp.ExitCode), nil // I can't pass around errors that well here, in channels, and THESE ERRORS OF COMPILATION AND ALL MUST BE SHOWN TO THE USER
	}

	return stdoutBuf.String(), nil
}

func (d *Docker) ExecInContainerStdCopy(command string) (string, error) {
	// Create an exec configuration
	execConfig := container.ExecOptions{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
	}

	// Create an exec instance in the container
	execID, err := d.cli.ContainerExecCreate(d.ctx, d.containerId, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %v", err)
	}

	// BIG RED FLAG, SOME CODE CHANGE HAS OCCURRED HERE
	resp, err := d.cli.ContainerExecAttach(d.ctx, execID.ID, container.ExecAttachOptions{Tty: false})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %v", err)
	}
	defer resp.Close()

	// Use buffers to capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer

	// Use stdcopy.StdCopy to separate stdout and stderr
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to copy output: %v", err)
	}

	// Combine stdout and stderr (if stderr is not empty)
	output := stdoutBuf.String()
	if stderrBuf.Len() > 0 {
		output += "\n" + stderrBuf.String()
	}

	return output, nil
}
