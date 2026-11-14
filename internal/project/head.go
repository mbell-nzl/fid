package project

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ObserveHead reads the committed head and whether anything needs committing.
// Unlike a full repository image, dirty paths need no content digest: hand-off
// refuses them, including unmerged paths. Every git child stays local and leaves
// the index alone. Hidden tracked changes are refused as in ObserveRepository.
func ObserveHead(root string) (Observation, error) {
	return captureRepository(func() (Observation, error) {
		observation := Observation{Root: root, WorktreeID: WorktreeID(root)}
		status, err := localHeadGit(root, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignore-submodules=none")
		if err != nil {
			return Observation{}, err
		}
		if len(status) != 0 {
			observation.Dirty = true
			return observation, nil
		}
		flags, err := localHeadGit(root, "ls-files", "-v", "-z")
		if err != nil {
			return Observation{}, err
		}
		for _, entry := range strings.Split(string(flags), "\x00") {
			if len(entry) > 0 && (entry[0] == 'S' || entry[0] >= 'a' && entry[0] <= 'z') {
				return Observation{}, unsafeRepositoryReason("tracked-flags", "tracked assume-unchanged or skip-worktree entry")
			}
		}
		head, err := localHeadGit(root, "rev-parse", "--verify", "HEAD^{commit}")
		if err != nil {
			return Observation{}, err
		}
		tree, err := localHeadGit(root, "rev-parse", "--verify", strings.TrimSpace(string(head))+"^{tree}")
		if err != nil {
			return Observation{}, err
		}
		observation.HeadOID, observation.TreeOID = strings.TrimSpace(string(head)), strings.TrimSpace(string(tree))
		return observation, nil
	})
}

func localHeadGit(root string, args ...string) ([]byte, error) {
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	command.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1", "GIT_OPTIONAL_LOCKS=0")
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("read hand-off %s: %w", args[0], err)
	}
	if err := checkRepositoryOutputSize(len(output)); err != nil {
		return nil, err
	}
	return output, nil
}
