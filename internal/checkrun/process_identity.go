package checkrun

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/checkrecord"
)

func machineIdentity() (string, error) {
	if runtime.GOOS == "linux" {
		if raw, err := os.ReadFile("/etc/machine-id"); err == nil {
			if id := strings.TrimSpace(string(raw)); id != "" {
				return "machine-id:" + id, nil
			}
		}
	}
	host, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("read machine: %w", err)
	}
	if host == "" {
		return "", fmt.Errorf("read machine: empty hostname")
	}
	return "hostname:" + host, nil
}

func processIdentity(pid int) (checkrecord.Process, error) {
	machine, err := machineIdentity()
	if err != nil {
		return checkrecord.Process{}, err
	}
	process := checkrecord.Process{Machine: machine, PID: pid}
	if runtime.GOOS == "linux" {
		boot, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
		if err != nil {
			return process, fmt.Errorf("read boot identity: %w", err)
		}
		process.BootID = strings.TrimSpace(string(boot))
		if process.BootID == "" {
			return process, fmt.Errorf("read boot identity: empty boot ID")
		}
		stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
		if err != nil {
			return process, fmt.Errorf("read process identity: %w", err)
		}
		index := strings.LastIndex(string(stat), ") ")
		fields, ok := "", index >= 0
		if ok {
			fields = string(stat)[index+2:]
		}
		values := strings.Fields(fields)
		if !ok || len(values) < 20 {
			return process, fmt.Errorf("invalid process identity")
		}
		process.StartTime = values[19]
		return process, nil
	}
	command := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "lstart=")
	output, err := command.Output()
	if err != nil || len(strings.TrimSpace(string(output))) == 0 {
		return process, fmt.Errorf("read process start time: %w", err)
	}
	process.StartTime = strings.TrimSpace(string(output))
	return process, nil
}

func processAlive(process checkrecord.Process) (bool, error) {
	if err := process.Validate(); err != nil {
		return false, err
	}
	self, err := processIdentity(os.Getpid())
	if err != nil {
		return false, err
	}
	if self.Machine != process.Machine {
		return false, fmt.Errorf("cannot inspect a check process on another machine")
	}
	if self.BootID == "" && process.BootID != "" {
		return false, fmt.Errorf("cannot inspect a check process with a different boot identity scheme")
	}
	if self.BootID != "" {
		if process.BootID == "" {
			return false, fmt.Errorf("cannot inspect a check process without its boot identity")
		}
		if self.BootID != process.BootID {
			return false, nil
		}
	}
	current, err := processIdentity(process.PID)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if runtime.GOOS == "linux" {
		stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", process.PID))
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		index := strings.LastIndex(string(stat), ") ")
		if index < 0 {
			return false, fmt.Errorf("invalid process state")
		}
		fields := strings.Fields(string(stat)[index+2:])
		if len(fields) == 0 {
			return false, fmt.Errorf("missing process state")
		}
		if fields[0] == "Z" || fields[0] == "X" {
			return false, nil
		}
	}
	return current == process, nil
}
