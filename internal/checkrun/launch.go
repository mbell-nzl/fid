package checkrun

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

const launchMarker = "__fid_check_child_v1"

type launchPipes struct{ permitReader, permitWriter, statusReader, statusWriter *os.File }

func launchCommand(executable string, argv []string) (*exec.Cmd, *launchPipes, error) {
	launcher, err := os.Executable()
	if err != nil {
		return nil, nil, fmt.Errorf("locate check launcher: %w", err)
	}
	pipes := &launchPipes{}
	pipes.permitReader, pipes.permitWriter, err = os.Pipe()
	if err != nil {
		return nil, nil, fmt.Errorf("create check permit pipe: %w", err)
	}
	pipes.statusReader, pipes.statusWriter, err = os.Pipe()
	if err != nil {
		pipes.close()
		return nil, nil, fmt.Errorf("create check status pipe: %w", err)
	}
	args := append([]string{launchMarker, executable}, argv...)
	command := exec.Command(launcher, args...)
	command.ExtraFiles = []*os.File{pipes.permitReader, pipes.statusWriter}
	return command, pipes, nil
}

func (pipes *launchPipes) closeChildEnds() {
	_ = pipes.permitReader.Close()
	_ = pipes.statusWriter.Close()
}

func (pipes *launchPipes) close() {
	for _, file := range []*os.File{pipes.permitReader, pipes.permitWriter, pipes.statusReader, pipes.statusWriter} {
		if file != nil {
			_ = file.Close()
		}
	}
}

// release permits exec only after the caller recorded the PID. EOF on the
// status pipe proves exec closed its private descriptor; data reports exec
// failure without confusing the command's own exit code with the launcher's.
func (pipes *launchPipes) release() (spawnFailed bool, err error) {
	if _, err := pipes.permitWriter.Write([]byte{1}); err != nil {
		return true, nil
	}
	if err := pipes.permitWriter.Close(); err != nil {
		return false, fmt.Errorf("close check permit: %w", err)
	}
	status, err := io.ReadAll(io.LimitReader(pipes.statusReader, 4096))
	_ = pipes.statusReader.Close()
	if err != nil {
		return false, fmt.Errorf("read check launch status: %w", err)
	}
	return len(status) > 0, nil
}

// LaunchChild is the private exec handshake used by the same fid executable.
// It is recognized only with both inherited pipe descriptors, never as a CLI
// leaf. A parent that dies before granting permission closes the pipe: the
// configured command cannot execute. Exec preserves the recorded PID and start
// identity and closes every private descriptor before entering the command.
func LaunchChild() (handled bool, err error) {
	if len(os.Args) < 4 || os.Args[1] != launchMarker {
		return false, nil
	}
	for _, fd := range []int{3, 4} {
		var stat syscall.Stat_t
		if err := syscall.Fstat(fd, &stat); err != nil || stat.Mode&syscall.S_IFMT != syscall.S_IFIFO {
			return false, nil
		}
	}
	permit, status := os.NewFile(3, "check-permit"), os.NewFile(4, "check-status")
	defer func() { _ = permit.Close(); _ = status.Close() }()
	var token [1]byte
	if _, err := io.ReadFull(permit, token[:]); err != nil || token[0] != 1 {
		return true, fmt.Errorf("check launch was not permitted")
	}
	if err := permit.Close(); err != nil {
		return true, fmt.Errorf("close check launch permit: %w", err)
	}
	syscall.CloseOnExec(4)
	err = syscall.Exec(os.Args[2], os.Args[3:], os.Environ())
	if err != nil {
		_, _ = io.WriteString(status, err.Error())
	}
	return true, err
}
