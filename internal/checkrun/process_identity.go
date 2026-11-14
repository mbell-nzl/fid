package checkrun

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/checkrecord"
)

func processIdentity(pid int) (checkrecord.Process, error) {
	host, err := os.Hostname()
	if err != nil {
		return checkrecord.Process{}, fmt.Errorf("read machine: %w", err)
	}
	process := checkrecord.Process{Machine: host, PID: pid}
	if runtime.GOOS == "linux" {
		boot, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
		if err != nil {
			return process, fmt.Errorf("read boot identity: %w", err)
		}
		process.Machine += ":" + strings.TrimSpace(string(boot))
		stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
		if err != nil {
			return process, fmt.Errorf("read process identity: %w", err)
		}
		_, fields, ok := strings.Cut(string(stat), ") ")
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

func processAlive(process checkrecord.Process) bool {
	current, err := processIdentity(process.PID)
	return err == nil && current == process
}
