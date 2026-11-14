package checkrun

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestProcessIdentitySeparatesStableMachineAndBoot(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "linux" {
		t.Skip("Linux boot identity")
	}
	process, err := processIdentity(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	machine, err := os.ReadFile("/etc/machine-id")
	if err == nil && strings.TrimSpace(string(machine)) != "" {
		if want := "machine-id:" + strings.TrimSpace(string(machine)); process.Machine != want {
			t.Fatalf("machine=%s, want stable %s", process.Machine, want)
		}
	}
	raw, err := json.Marshal(process)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	boot, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
	if err != nil {
		t.Fatal(err)
	}
	if fields["boot_id"] != strings.TrimSpace(string(boot)) {
		t.Fatalf("separate boot identity missing: %s", raw)
	}
}

func TestProcessAliveRequiresTheSameMachineAndBoot(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "linux" {
		t.Skip("Linux boot identity")
	}
	for _, row := range []struct {
		name, machine, boot string
		alive, refused      bool
	}{
		{name: "current process", alive: true},
		{name: "prior boot", boot: "00000000-0000-0000-0000-000000000000"},
		{name: "foreign machine", machine: "machine-id:another", refused: true},
		{name: "foreign machine prior boot", machine: "machine-id:another", boot: "00000000-0000-0000-0000-000000000000", refused: true},
		{name: "uninspectable boot", boot: "unknown", refused: true},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			process, err := processIdentity(os.Getpid())
			if err != nil {
				t.Fatal(err)
			}
			raw, err := json.Marshal(process)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]any
			if err := json.Unmarshal(raw, &fields); err != nil {
				t.Fatal(err)
			}
			if row.machine != "" {
				fields["machine"] = row.machine
			}
			if row.boot != "" {
				fields["boot_id"] = row.boot
			}
			raw, err = json.Marshal(fields)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(raw, &process); err != nil {
				t.Fatal(err)
			}
			alive, err := processAlive(process)
			if alive != row.alive || (err != nil) != row.refused {
				t.Fatalf("liveness=%v/%v, want %v/refused=%v", alive, err, row.alive, row.refused)
			}
		})
	}
}
