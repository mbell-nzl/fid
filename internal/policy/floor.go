package policy

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/pathglob"
)

// Floor explains the highest lane required by committed changed paths. Nil
// lanes mean no floor; only built-in hard paths prevent a person's waiver.
type Floor struct {
	Lane      *string  `json:"lane"`
	Paths     []string `json:"paths"`
	HardFloor *string  `json:"hard_floor"`
	HardPaths []string `json:"hard_paths"`
}

// Raise preserves a higher planned lane. An empty change never lowers it.
func (floor Floor) Raise(lane string) string {
	if floor.Lane != nil && laneRank(*floor.Lane) > laneRank(lane) {
		return *floor.Lane
	}
	return lane
}

func laneRank(lane string) int {
	switch lane {
	case "routine":
		return 1
	case "standard":
		return 2
	case "assured":
		return 3
	}
	return 0
}

// ComputeFloor is shared by checking, hand-off admission and display. Paths
// come from Compare, which includes deletions and both names of renames.
func ComputeFloor(effective Effective, paths []string) (Floor, error) {
	floor := Floor{Paths: []string{}, HardPaths: []string{}}
	hard, err := compileFloorGlobs(effective.HardAssured)
	if err != nil {
		return floor, err
	}
	assuredPatterns := append([]string{}, effective.Lanes["assured"].Paths...)
	for _, check := range effective.Checks {
		assuredPatterns = append(assuredPatterns, check.DefinedBy...)
	}
	assured, err := compileFloorGlobs(assuredPatterns)
	if err != nil {
		return floor, err
	}
	routine, err := compileFloorGlobs(effective.Lanes["routine"].Paths)
	if err != nil {
		return floor, err
	}
	for _, path := range canonicalSet(paths) {
		lane := "standard"
		isHard := matchesFloor(hard, path)
		switch {
		case isHard || !pathglob.ValidPath(path) || matchesFloor(assured, path):
			lane = "assured"
		case matchesFloor(routine, path):
			lane = "routine"
		}
		if isHard {
			value := "assured"
			floor.HardFloor = &value
			floor.HardPaths = append(floor.HardPaths, path)
		}
		if floor.Lane == nil || laneRank(lane) > laneRank(*floor.Lane) {
			floor.Lane = &lane
			floor.Paths = []string{}
		}
		if lane == *floor.Lane {
			floor.Paths = append(floor.Paths, string([]rune(path)))
		}
	}
	return floor, nil
}

func compileFloorGlobs(patterns []string) ([]pathglob.Glob, error) {
	globs := []pathglob.Glob{}
	for _, pattern := range patterns {
		glob, err := pathglob.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("lane floor: %w", err)
		}
		globs = append(globs, glob)
	}
	return globs, nil
}

func matchesFloor(globs []pathglob.Glob, path string) bool {
	for _, glob := range globs {
		if glob.Match(path) {
			return true
		}
	}
	return false
}

// PathBytes preserves Git path bytes through JSON (as base64), including names
// which JSON strings cannot represent without replacing invalid UTF-8.
func PathBytes(paths []string) [][]byte {
	result := make([][]byte, len(paths))
	for i, path := range paths {
		result[i] = []byte(path)
	}
	return result
}

func PathsFromBytes(paths [][]byte) []string {
	result := make([]string, len(paths))
	for i, path := range paths {
		result[i] = string(path)
	}
	return result
}
