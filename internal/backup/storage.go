package backup

import (
	"archive/zip"
	"fmt"
	"strconv"
	"strings"
)

// packClone bounds the loose history left by plumbing commands, which do
// not trigger git's automatic maintenance. Count exactly instead of git's
// gc.auto sampling, and run in the foreground while holding the clone lock.
// Before export, every successful backup is already protected by its ref.
// Repack keeps unreachable objects; prune-packed only removes loose copies
// of objects already in a pack, so failed pushes are never pruned early.
func packClone(repository clone) error {
	count, err := repository.text("count-objects")
	if err != nil {
		return failure(KindInternal, err, "count the backup clone's loose objects")
	}
	loose, _, _ := strings.Cut(count, " ")
	number, err := strconv.Atoi(loose)
	if err != nil {
		return failure(KindInternal, err, "read the backup clone's loose object count")
	}
	if number < 256 {
		return nil
	}
	for _, command := range [][]string{{"repack", "-A", "-d"}, {"prune-packed"}} {
		if _, err := repository.run(nil, command...); err != nil {
			return failure(KindInternal, err, "pack the backup clone")
		}
	}
	return nil
}

type memberTree struct {
	blobs       map[string]string
	directories map[string]*memberTree
}

// storeMembers streams the validated export's members into git's object
// store without extracting them to disk. Stdin has no path, so attributes
// and filters from the user's config cannot change the contents.
func storeMembers(repository clone, archive, manifestBlob string) (string, error) {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return "", err
	}
	defer func() { _ = reader.Close() }()
	root := &memberTree{blobs: map[string]string{ManifestName: manifestBlob}}
	for _, member := range reader.File {
		input, err := member.Open()
		if err != nil {
			return "", err
		}
		blob, err := repository.output(input, "hash-object", "-w", "--stdin")
		closeErr := input.Close()
		if err != nil {
			return "", err
		}
		if closeErr != nil {
			return "", closeErr
		}
		root.add(member.Name, blob)
	}
	return root.write(repository)
}

func (tree *memberTree) add(path, blob string) {
	name, rest, directory := strings.Cut(path, "/")
	if !directory {
		if tree.blobs == nil {
			tree.blobs = map[string]string{}
		}
		tree.blobs[name] = blob
		return
	}
	if tree.directories == nil {
		tree.directories = map[string]*memberTree{}
	}
	if tree.directories[name] == nil {
		tree.directories[name] = &memberTree{}
	}
	tree.directories[name].add(rest, blob)
}

func (tree *memberTree) write(repository clone) (string, error) {
	var entries strings.Builder
	for name, blob := range tree.blobs {
		fmt.Fprintf(&entries, "100644 blob %s\t%s\x00", blob, name)
	}
	for name, directory := range tree.directories {
		oid, err := directory.write(repository)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&entries, "040000 tree %s\t%s\x00", oid, name)
	}
	return repository.output(strings.NewReader(entries.String()), "mktree", "-z")
}
