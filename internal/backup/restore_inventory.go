package backup

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/export"
)

var errFolderResourceLimit = errors.New("archive resource limit exceeded")

type folderMember struct {
	name string
	oid  string
}

// folderInventory consumes Git's NUL-delimited metadata incrementally. The
// entry ceiling applies before the listing or member inventory can grow
// without bound; sizes are checked before reading or recompressing any blob.
type folderInventory struct {
	id      string
	limits  export.ResourceLimits
	pending bytes.Buffer
	members []folderMember
	total   int64
	failed  error
}

func listFolderMembers(repository clone, commit, id string, limits export.ResourceLimits) ([]folderMember, error) {
	inventory := &folderInventory{id: id, limits: limits}
	err := repository.stream(inventory, "ls-tree", "-r", "-l", "-z", commit, "--", id+"/")
	if inventory.failed != nil {
		return nil, inventory.failed
	}
	if err != nil {
		return nil, failure(KindInternal, err, "read the backup folder's members")
	}
	return inventory.members, nil
}

func (inventory *folderInventory) Write(payload []byte) (int, error) {
	count := len(payload)
	for len(payload) != 0 {
		end := bytes.IndexByte(payload, 0)
		if end < 0 {
			_, _ = inventory.pending.Write(payload)
			break
		}
		_, _ = inventory.pending.Write(payload[:end])
		inventory.failed = inventory.add(inventory.pending.String())
		if inventory.failed != nil {
			return count, inventory.failed
		}
		inventory.pending.Reset()
		payload = payload[end+1:]
	}
	return count, nil
}

func (inventory *folderInventory) add(record string) error {
	metadata, path, _ := strings.Cut(record, "\t")
	if path == inventory.id+"/"+ManifestName {
		return nil
	}
	fields := strings.Fields(metadata)
	if len(fields) != 4 || fields[0] != "100644" || fields[1] != "blob" {
		return damaged(inventory.id, folderCheck(CheckManifestInvalid, "%s is not a regular export member", path))
	}
	name := strings.TrimPrefix(path, inventory.id+"/")
	size, err := strconv.ParseInt(fields[3], 10, 64)
	limit := inventory.limits.NonArtifactBytes
	if strings.HasPrefix(name, "artifacts/") {
		limit = inventory.limits.ArtifactBytes
	}
	if err != nil || size < 0 || size > limit || size > inventory.limits.ExpandedBytes-inventory.total || len(inventory.members) >= inventory.limits.Entries {
		return folderResourceLimit(inventory.id, path)
	}
	inventory.total += size
	inventory.members = append(inventory.members, folderMember{name: name, oid: fields[2]})
	return nil
}

func folderResourceLimit(id, path string) error {
	return damaged(id, folderCheck("resource-limit", "%s exceeds the export's resource limits", path))
}

// Size metadata rejects large manifests before allocation. Keep the actual
// write bounded too, so no Git output can grow either manifest buffer past
// the same non-Artifact ceiling.
type folderBlobWriter struct {
	output    io.Writer
	remaining int64
	exceeded  bool
}

func (writer *folderBlobWriter) Write(payload []byte) (int, error) {
	if int64(len(payload)) > writer.remaining {
		writer.exceeded = true
		return 0, errFolderResourceLimit
	}
	n, err := writer.output.Write(payload)
	writer.remaining -= int64(n)
	return n, err
}
