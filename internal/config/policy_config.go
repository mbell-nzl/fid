package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mbell-nzl/fid/internal/pathglob"
)

// ProjectFile is the shared identity and policy schema of .fid.toml.
// Identity consumers require ProjectID; policy consumers allow it to be absent.
type ProjectFile struct {
	ProjectID       string
	Checks          map[string]Check
	Lanes           map[string]Lane
	TestAdequacy    []string
	TestAdequacySet bool
}

type Check struct {
	Run         []string      `json:"run"`
	Timeout     time.Duration `json:"timeout_ns"`
	DefinedBy   []string      `json:"defined_by"`
	WhenChanged []string      `json:"when_changed"`
}

type Lane struct {
	Paths       []string `json:"paths"`
	Checks      []string `json:"checks"`
	ChecksSet   bool     `json:"checks_set"`
	ReworkCap   int64    `json:"rework_cap"`
	TokenBudget int64    `json:"token_budget"`
	AutoMerge   bool     `json:"auto_merge"`
}

type projectDocument struct {
	Project struct {
		ID string `toml:"id"`
	} `toml:"project"`
	Checks map[string]struct {
		Run         []string `toml:"run"`
		Timeout     string   `toml:"timeout"`
		DefinedBy   []string `toml:"defined_by"`
		WhenChanged []string `toml:"when_changed"`
	} `toml:"checks"`
	Lanes map[string]struct {
		Paths       []string `toml:"paths"`
		Checks      []string `toml:"checks"`
		ReworkCap   *int64   `toml:"rework_cap"`
		TokenBudget int64    `toml:"token_budget"`
		AutoMerge   bool     `toml:"auto_merge"`
	} `toml:"lanes"`
	Review struct {
		TestAdequacy []string `toml:"test_adequacy"`
	} `toml:"review"`
}

// DecodeProject refuses unknown keys and validates the complete file,
// including policy when the caller only needs its project identity.
func DecodeProject(data []byte) (ProjectFile, error) {
	var document projectDocument
	metadata, err := toml.Decode(string(data), &document)
	if err != nil {
		return ProjectFile{}, fmt.Errorf("parse .fid.toml: %w", err)
	}
	if unknown := metadata.Undecoded(); len(unknown) != 0 {
		return ProjectFile{}, fmt.Errorf("parse .fid.toml: unknown key %s", unknown[0])
	}
	// Struct decoding is case-insensitive, but TOML keys and IsDefined are
	// case-sensitive. Check every original key, including inline/dotted keys,
	// before using decoded values or presence checks.
	for _, key := range metadata.Keys() {
		if !projectKeyKnown(key) {
			return ProjectFile{}, fmt.Errorf("parse .fid.toml: unknown key %s", key)
		}
	}
	if metadata.IsDefined("project") && !validBackupProjectID(document.Project.ID) {
		return ProjectFile{}, fmt.Errorf("project.id must be project- followed by 32 lowercase hexadecimal digits")
	}
	result := ProjectFile{ProjectID: document.Project.ID, Checks: map[string]Check{}, Lanes: map[string]Lane{}, TestAdequacy: document.Review.TestAdequacy, TestAdequacySet: metadata.IsDefined("review", "test_adequacy")}
	if err := decodeChecks(document, &result); err != nil {
		return ProjectFile{}, err
	}
	if err := decodeLanes(document, metadata, &result); err != nil {
		return ProjectFile{}, err
	}
	if err := validateGlobs("review.test_adequacy", result.TestAdequacy); err != nil {
		return ProjectFile{}, err
	}
	return result, nil
}

func projectKeyKnown(key toml.Key) bool {
	switch key[0] {
	case "project":
		return len(key) == 1 || len(key) == 2 && key[1] == "id"
	case "checks":
		if len(key) <= 2 { // Check names are project-defined and case-sensitive.
			return true
		}
		if len(key) == 3 {
			switch key[2] {
			case "run", "timeout", "defined_by", "when_changed":
				return true
			}
		}
	case "lanes":
		if len(key) == 1 {
			return true
		}
		if key[1] != "routine" && key[1] != "standard" && key[1] != "assured" {
			return false
		}
		if len(key) == 2 {
			return true
		}
		if len(key) == 3 {
			switch key[2] {
			case "paths", "checks", "rework_cap", "token_budget", "auto_merge":
				return true
			}
		}
	case "review":
		return len(key) == 1 || len(key) == 2 && key[1] == "test_adequacy"
	}
	return false
}

func decodeChecks(document projectDocument, result *ProjectFile) error {
	for name, check := range document.Checks {
		if len(check.Run) == 0 || check.Run[0] == "" {
			return fmt.Errorf("checks.%s.run must be non-empty argv with an executable", name)
		}
		timeout, err := time.ParseDuration(check.Timeout)
		if err != nil || timeout <= 0 {
			return fmt.Errorf("checks.%s.timeout must be a positive Go duration", name)
		}
		for field, globs := range map[string][]string{"defined_by": check.DefinedBy, "when_changed": check.WhenChanged} {
			if err := validateGlobs("checks."+name+"."+field, globs); err != nil {
				return err
			}
		}
		result.Checks[name] = Check{check.Run, timeout, check.DefinedBy, check.WhenChanged}
	}
	return nil
}

func decodeLanes(document projectDocument, metadata toml.MetaData, result *ProjectFile) error {
	for _, name := range []string{"routine", "standard", "assured"} {
		result.Lanes[name] = Lane{ReworkCap: 3}
	}
	for name, lane := range document.Lanes {
		if _, found := result.Lanes[name]; !found {
			return fmt.Errorf("unknown lane %q", name)
		}
		if name == "standard" && metadata.IsDefined("lanes", name, "paths") {
			return fmt.Errorf("lanes.standard.paths is refused: other paths are Standard")
		}
		if name == "assured" && metadata.IsDefined("lanes", name, "auto_merge") {
			return fmt.Errorf("lanes.assured.auto_merge is refused: a person always merges Assured work")
		}
		cap := int64(3)
		if lane.ReworkCap != nil {
			cap = *lane.ReworkCap
		}
		if cap < 0 || lane.TokenBudget < 0 {
			return fmt.Errorf("lanes.%s.rework_cap and token_budget must be non-negative integers", name)
		}
		for _, check := range lane.Checks {
			if _, found := result.Checks[check]; !found {
				return fmt.Errorf("lanes.%s.checks names undefined check %q", name, check)
			}
		}
		if err := validateGlobs("lanes."+name+".paths", lane.Paths); err != nil {
			return err
		}
		result.Lanes[name] = Lane{lane.Paths, lane.Checks, metadata.IsDefined("lanes", name, "checks"), cap, lane.TokenBudget, lane.AutoMerge}
	}
	return nil
}

func validateGlobs(field string, globs []string) error {
	for _, pattern := range globs {
		if _, err := pathglob.Compile(pattern); err != nil {
			return fmt.Errorf("%s: %w", field, err)
		}
	}
	return nil
}
