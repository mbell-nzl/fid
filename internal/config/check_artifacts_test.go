package config

import (
	"reflect"
	"testing"
)

func TestCheckArtifactsArePlainPaths(t *testing.T) {
	t.Parallel()
	file, err := DecodeProject([]byte("[checks.test]\nrun = ['true']\ntimeout = '15m'\nartifacts = ['coverage.out', 'reports/a[1].txt']\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(file.Checks["test"].Artifacts, []string{"coverage.out", "reports/a[1].txt"}) {
		t.Fatalf("artifacts = %#v", file.Checks["test"].Artifacts)
	}
}
